use crate::com::iabtechlab::adcom::v1::enums::OperatingSystem;
use crate::com::iabtechlab::adcom::v1 as adcom;
use crate::com::iabtechlab::adcom::v1::context::DistributionChannel;
use crate::com::iabtechlab::openrtb::v3 as openrtb;
use crate::com::iabtechlab::openrtb::v3::AuctionType;
use crate::sdk;
use crate::sdk::{AdFormat, AuctionRequest, DeviceConnectionType, DeviceType, Geo, Segment};
use crate::org::bidon::proto::v1::mediation;
use crate::org::bidon::proto::v1::mediation::{
    Demand, DeviceExt, Orientation, APP_EXT, AUCTION_RESPONSE_EXT, BID_EXT, DEVICE_EXT,
    PLACEMENT_EXT, REGS_EXT, USER_EXT,
};
use anyhow::{anyhow, Result};
use sdk::AdObjectOrientation;
use prost::{Extendable, Message};
use serde_json::Value;
use std::collections::HashMap;

//TODO As it takes auction_request's ownership, it should be possible to remove most of the .clone() calls.
pub(crate) fn try_from(
    auction_request: AuctionRequest,
    bidon_version: &String,
) -> Result<openrtb::Openrtb> {
    // Convert AuctionRequest to Openrtb::Request
    let request = openrtb::Request {
        id: auction_request.ad_object.auction_id.to_owned(),
        test: auction_request.test,
        tmax: auction_request.tmax.map(|t| t as u32),
        at: Some(AuctionType::FirstPrice as i32),
        context: Some(serialize_context(&auction_request, bidon_version)?),
        item: vec![convert_ad_object_to_item(&auction_request.ad_object)?],
        ..Default::default()
    };

    // Create Openrtb instance with the converted request
    Ok(openrtb::Openrtb {
        ver: Some("3.0".to_string()),
        domainspec: Some("domain_spec".to_string()),
        domainver: Some("domain_version".to_string()),
        payload_oneof: Some(openrtb::openrtb::PayloadOneof::Request(request)),
    })
}

fn serialize_context(auction_request: &AuctionRequest, bidon_version: &String) -> Result<Vec<u8>> {
    // Create the AdCOM Context message
    let context = crate::org::bidon::proto::v1::context::Context {
        distribution_channel: DistributionChannel {
            channel_oneof: Some(adcom::context::distribution_channel::ChannelOneof::App(
                convert_app(&auction_request.app, bidon_version)?,
            )),
            ..Default::default()
        }
        .into(),
        device: convert_device(&auction_request)?.into(),
        user: convert_user(&auction_request.user, auction_request.segment.as_ref())?.into(),
        regs: match auction_request.regs.as_ref() {
            Some(t) => convert_regs(&t)?.into(),
            None => None,
        },
        restrictions: None, // TODO
    };

    // Serialize the Context message into bytes
    let mut context_bytes = Vec::new();
    context.encode(&mut context_bytes)?;
    Ok(context_bytes)
}

fn convert_app(
    api_app: &sdk::App,
    bidon_version: &String,
) -> Result<adcom::context::distribution_channel::App> {
    let mut app = adcom::context::distribution_channel::App {
        ver: api_app.version.clone().into(),
        keywords: None,
        paid: None,
        bundle: api_app.bundle.clone().into(),
        ..Default::default()
    };

    let bidon_app = crate::org::bidon::proto::v1::mediation::AppExt {
        key: api_app.key.clone().into(),
        framework: api_app.framework_version.clone(),
        framework_version: api_app.framework_version.clone(),
        plugin_version: api_app.plugin_version.clone(),
        sdk_version: api_app.sdk_version.clone(),
        skadn: api_app.skadn.clone().unwrap_or(vec![]),
        bidon_version: Some(bidon_version.clone()),
    };

    app.set_extension_data(APP_EXT, bidon_app)?;
    Ok(app)
}

fn convert_device(request: &AuctionRequest) -> Result<adcom::context::Device> {
    let api_device: &sdk::Device = &request.device;
    let geo: Option<&Geo> = request.geo.as_ref();

    let mut device = adcom::context::Device {
        // Map standard fields
        r#type: convert_device_type(api_device.device_type.clone()).map(|dt| dt as i32),
        ua: api_device.ua.clone().into(),
        make: api_device.make.clone().into(),
        model: api_device.model.clone().into(),
        os: Some(convert_os(api_device.os.clone()) as i32),
        osv: api_device.osv.clone().into(),
        hwv: api_device.hwv.clone().into(),
        h: api_device.h.into(),
        w: api_device.w.into(),
        ppi: api_device.ppi.into(),
        pxratio: (api_device.pxratio as f32).into(), // TODO validate conversion
        js: Some(api_device.js != 0),
        lang: api_device.language.clone().into(),
        carrier: api_device.clone().carrier,
        mccmnc: api_device.clone().mccmnc,
        contype: Some(<adcom::enums::ConnectionType as Into<i32>>::into(
            convert_connection_type(api_device.connection_type),
        )),
        geo: geo.clone().map(|g| convert_geo(&g)),
        ..Default::default()
    };

    let bidon_device_ext = convert_session(&request.session);
    device.set_extension_data(DEVICE_EXT, bidon_device_ext)?;

    Ok(device)
}

fn convert_device_type(device_type: Option<DeviceType>) -> Option<adcom::enums::DeviceType> {
    device_type
        .map(|dt| match dt {
            DeviceType::Phone => adcom::enums::DeviceType::Phone,
            DeviceType::Tablet => adcom::enums::DeviceType::Tablet,
        })
        .map(Into::into)
}

fn convert_os(os: String) -> OperatingSystem {
    match os.to_lowercase().as_str() {
        "ios" => OperatingSystem::Ios,
        "android" => OperatingSystem::Android,
        "windows" => OperatingSystem::Windows,
        "macos" => OperatingSystem::Macos,
        "linux" => OperatingSystem::Linux,
        _ => OperatingSystem::OtherNotListed,
    }
}

fn convert_connection_type(connection_type: DeviceConnectionType) -> adcom::enums::ConnectionType {
    match connection_type {
        DeviceConnectionType::Ethernet => adcom::enums::ConnectionType::Wired,
        DeviceConnectionType::Wifi => adcom::enums::ConnectionType::Wifi,
        DeviceConnectionType::CellularUnknown => adcom::enums::ConnectionType::CellUnknown,
        DeviceConnectionType::Cellular => adcom::enums::ConnectionType::CellUnknown,
        DeviceConnectionType::Cellular2G => adcom::enums::ConnectionType::Cell2g,
        DeviceConnectionType::Cellular3G => adcom::enums::ConnectionType::Cell3g,
        DeviceConnectionType::Cellular4G => adcom::enums::ConnectionType::Cell4g,
        DeviceConnectionType::Cellular5G => adcom::enums::ConnectionType::Cell5g,
    }
}

fn convert_geo(geo: &sdk::Geo) -> adcom::context::Geo {
    let geo = adcom::context::Geo {
        r#type: Some(adcom::enums::LocationType::Unknown as i32), // TODO
        lat: geo.lat.map(|t| t as f32),
        lon: geo.lon.map(|t| t as f32),
        accur: geo.accuracy.map(|t| (t as i32)), // TODO check accuracy conversion. We convert it from f64 to i32 here.
        country: geo.country.clone().into(),
        city: geo.city.clone().into(),
        zip: geo.zip.clone().into(),
        utcoffset: geo.utcoffset,
        lastfix: geo.lastfix,
        ..Default::default()
    };
    geo
}

fn convert_user(
    api_user: &sdk::User,
    segment: Option<&Segment>,
) -> Result<adcom::context::User> {
    let mut user = adcom::context::User {
        id: api_user.idg.map(|uuid| uuid.to_string()),
        consent: api_user
            .consent
            .as_ref()
            .and_then(|c| serde_json::to_string(c).ok()),
        ..Default::default()
    };

    let bidon_user_ext = crate::org::bidon::proto::v1::mediation::UserExt {
        idfa: api_user.idfa.map(|uuid| uuid.to_string()),
        tracking_authorization_status: Some(api_user.tracking_authorization_status.clone()),
        idfv: api_user.idfv.map(|uuid| uuid.to_string()),
        idg: api_user.idg.map(|uuid| uuid.to_string()),
        segments: segment.into_iter().map(convert_segment).collect(),
    };

    user.set_extension_data(USER_EXT, bidon_user_ext)?;

    Ok(user)
}

fn convert_segment(api_segment: &Segment) -> crate::org::bidon::proto::v1::mediation::Segment {
    let segment = crate::org::bidon::proto::v1::mediation::Segment {
        id: api_segment.id.clone(),
        uid: api_segment.uid.clone(),
        ext: api_segment.ext.clone(),
    };

    segment
}

fn convert_session(api_session: &sdk::Session) -> DeviceExt {
    DeviceExt {
        id: Some(api_session.id.to_string().clone()),
        launch_ts: Some(api_session.launch_ts),
        launch_monotonic_ts: api_session.launch_monotonic_ts.into(),
        start_ts: Some(api_session.start_ts),
        start_monotonic_ts: Some(api_session.start_monotonic_ts),
        ts: Some(api_session.ts),
        monotonic_ts: Some(api_session.monotonic_ts),
        memory_warnings_ts: api_session.memory_warnings_ts.clone(),
        memory_warnings_monotonic_ts: api_session.memory_warnings_monotonic_ts.clone(),
        ram_used: Some(api_session.ram_used),
        ram_size: Some(api_session.ram_size),
        storage_free: api_session.storage_free,
        storage_used: api_session.storage_used,
        battery: Some(api_session.battery),
        cpu_usage: Some(api_session.cpu_usage),
    }
}

fn convert_regs(api_regs: &sdk::Regulations) -> Result<adcom::context::Regs> {
    let mut regs = adcom::context::Regs {
        coppa: api_regs.coppa,
        gdpr: api_regs.gdpr.clone(),
        ..Default::default()
    };

    let mediation_regs = crate::org::bidon::proto::v1::mediation::RegsExt {
        us_privacy: api_regs.us_privacy.clone(),
        eu_privacy: api_regs.eu_privacy.clone(),
        iab: match api_regs.iab.as_ref() {
            Some(t) => Some(convert_iab(t)?),
            None => None,
        },
    };

    regs.set_extension_data(REGS_EXT, mediation_regs)?;

    Ok(regs)
}

fn convert_iab(iab_json: &HashMap<String, Value>) -> Result<String> {
    serde_json::to_string(&iab_json).map_err(Into::into)
}

fn create_placement(ad_object: &sdk::AdObject) -> Result<adcom::placement::Placement> {
    let mut placement = crate::com::iabtechlab::adcom::v1::placement::Placement {
        display: None,
        video: None,
        audio: None, // Assuming no audio placement in this context
        // Common placement properties
        secure: Some(1), // Assuming HTTPS is required
        ..Default::default()
    };

    // Create the AdObjectExtension
    let placement_ext = crate::org::bidon::proto::v1::mediation::PlacementExt {
        auction_id: ad_object.auction_id.clone(),
        auction_key: ad_object.auction_key.clone(),
        auction_configuration_id: ad_object.auction_configuration_id.clone(),
        auction_configuration_uid: ad_object.auction_configuration_uid.clone(),
        orientation: ad_object
            .orientation
            .as_ref()
            .map(convert_ad_orientation)
            .map(|f| f as i32),
        demands: convert_demand(&ad_object.demands)?,
        banner: match &ad_object.banner {
            Some(ref banner) => Some(convert_banner_ad(banner)),
            None => None,
        },
        interstitial: ad_object.interstitial.as_ref().map(|i| i.to_string()),
        rewarded: ad_object.rewarded.as_ref().map(|r| r.to_string()),
    };

    placement.set_extension_data(PLACEMENT_EXT, placement_ext)?;
    Ok(placement)
}

fn serialize_placement(placement: &adcom::placement::Placement) -> Result<Vec<u8>> {
    let mut bytes = Vec::new();
    placement.encode(&mut bytes)?;
    Ok(bytes)
}

fn convert_ad_object_to_item(ad_object: &sdk::AdObject) -> Result<openrtb::Item> {
    let item = openrtb::Item {
        id: ad_object.auction_id.clone(),
        flr: Some(ad_object.auction_pricefloor as f32),
        flrcur: Some("USD".to_string()),
        spec: Some(serialize_placement(&create_placement(ad_object)?)?),
        ..Default::default()
    };

    Ok(item)
}

fn convert_ad_orientation(orientation: &AdObjectOrientation) -> Orientation {
    match orientation {
        AdObjectOrientation::Portrait => Orientation::Portrait,
        AdObjectOrientation::Landscape => Orientation::Landscape,
    }
}

fn convert_banner_ad(api_banner: &sdk::BannerAdObject) -> mediation::BannerAd {
    let banner = mediation::BannerAd {
        format: Some(convert_ad_format(api_banner.format) as i32),
    };

    banner
}

fn convert_ad_format(format: AdFormat) -> mediation::AdFormat {
    match format {
        AdFormat::Banner => mediation::AdFormat::Banner,
        AdFormat::Leaderboard => mediation::AdFormat::Leaderboard,
        AdFormat::Mrec => mediation::AdFormat::Mrec,
        AdFormat::Adaptive => mediation::AdFormat::Adaptive,
    }
}

fn convert_demand(api_demand: &HashMap<String, Value>) -> Result<HashMap<String, Demand>> {
    let mut demands = HashMap::new();

    for (key, value) in api_demand {
        let map = value
            .as_object()
            .ok_or(anyhow!("Demand value is not an object: {}", value))?;
        let demand = Demand {
            // Assuming Demand has fields that need to be populated from the value
            // Add the necessary field mappings here
            token: match map.get("token") {
                Some(v) => Some(
                    v.as_str()
                        .ok_or(anyhow!(
                            "Token is not a string. Key: {}, value: {}",
                            key,
                            value
                        ))?
                        .to_string(),
                ),
                None => None,
            },
            status: match map.get("status") {
                Some(v) => Some(
                    v.as_str()
                        .ok_or(anyhow!(
                            "Status is not a string. Key: {}, value: {}",
                            key,
                            value
                        ))?
                        .to_string(),
                ),
                None => None,
            },
            token_finish_ts: match map.get("token_finish_ts") {
                Some(v) => Some(v.as_i64().ok_or(anyhow!(
                    "token_finish_ts is not a number. Key: {}, value: {}",
                    key,
                    value
                ))?),
                None => None,
            },
            token_start_ts: match map.get("token_start_ts") {
                Some(v) => Some(v.as_i64().ok_or(anyhow!(
                    "token_start_ts is not a number. Key: {}, value: {}",
                    key,
                    value
                ))?),
                None => None,
            },
        };
        demands.insert(key.clone(), demand);
    }
    Ok(demands)
}

// TODO: this is a temporary function to convert the OpenRTB response to the AuctionResponse
// fix it later
pub(crate) fn try_into(openrtb: openrtb::Openrtb) -> Result<sdk::AuctionResponse> {
    // Extract the Response from Openrtb
    let response = match openrtb.payload_oneof {
        Some(openrtb::openrtb::PayloadOneof::Response(response)) => response,
        _ => return Err(anyhow!("OpenRTB payload is not a Response")),
    };

    // Extract bid information from the response
    let mut ad_units = Vec::new();
    let mut no_bids = Vec::new();

    for seatbid in response.seatbid {
        for bid in seatbid.bid {
            // Extract bid extension data
            let bid_ext = bid
                .extension_set
                .extension_data(BID_EXT)
                .map_err(|_| anyhow!("Missing mediation ad object extension in bid"))?;
            let ad_unit = sdk::AdUnit {
                label: bid_ext.label.clone().ok_or(anyhow!("Label is missing"))?,
                uid: bid.item.unwrap_or_default(),
                demand_id: bid.cid.unwrap_or_default(),
                // Start of Selection
                pricefloor: Some(bid.price.unwrap_or_default() as f64),
                bid_type: bid_ext
                    .bid_type
                    .clone()
                    .ok_or(anyhow!("Bid type is missing"))?,
                ext: Some(
                    serde_json::to_value(bid_ext.ext.clone())
                        .unwrap_or_default()
                        .as_object()
                        .map(|map| map.clone().into_iter().collect::<HashMap<String, Value>>())
                        .unwrap_or_default(),
                ),
            };
            if bid.price.unwrap_or_default() > 0.0 {
                ad_units.push(ad_unit);
            } else {
                no_bids.push(ad_unit);
            }
        }
    }

    // Extract auction configuration from response extensions
    let auction_ext = response
        .extension_set
        .extension_data(AUCTION_RESPONSE_EXT)
        .map_err(|_| anyhow!("Missing mediation ad object extension in response"))?;

    let auction_response = sdk::AuctionResponse {
        ad_units,
        auction_id: response.id.unwrap_or_default(),
        no_bids: Some(no_bids),
        token: auction_ext.token.clone().unwrap_or_default(),
        external_win_notifications: auction_ext
            .external_win_notifications
            .clone()
            .unwrap_or_default(),
        segment: auction_ext
            .segment
            .as_ref()
            .map(|s| sdk::Segment {
                id: s.id.clone(),
                uid: s.uid.clone(),
                ext: s.ext.clone(),
            })
            .ok_or(anyhow!("Segment is missing"))?,
        auction_configuration_id: auction_ext.auction_configuration_id.unwrap_or_default(),
        auction_configuration_uid: auction_ext
            .auction_configuration_uid
            .clone()
            .unwrap_or_default(),
        auction_pricefloor: auction_ext.auction_pricefloor.unwrap_or_default(),
        auction_timeout: auction_ext.auction_timeout.unwrap_or_default(),
    };

    Ok(auction_response)
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::sdk::{App, Session};
    use crate::org::bidon::proto::v1::mediation::USER_EXT;
    use prost::{Extension, ExtensionRegistry};
    use serde_json::json;
    use std::collections::HashMap;
    use uuid::Uuid;

    fn create_test_auction_request() -> sdk::AuctionRequest {
        sdk::AuctionRequest {
            ad_object: sdk::AdObject {
                auction_id: Some("auction123".to_string()),
                auction_key: Some("key123".to_string()),
                auction_configuration_id: Some(456),
                auction_configuration_uid: Some("config789".to_string()),
                auction_pricefloor: 1.0,
                orientation: None,
                demands: HashMap::new(),
                banner: None,
                interstitial: None,
                rewarded: None,
            },
            adapters: HashMap::new(),
            app: App {
                bundle: "com.example.app".to_string(),
                framework: "".to_string(),
                framework_version: None,
                key: "".to_string(),
                plugin_version: None,
                sdk_version: None,
                skadn: None,
                version: "".to_string(),
            },
            device: sdk::Device {
                device_type: Some(DeviceType::Phone),
                ua: "Mozilla/5.0".to_string(),
                make: "Apple".to_string(),
                model: "iPhone".to_string(),
                os: "iOS".to_string(),
                osv: "14.4".to_string(),
                hwv: "A14".to_string(),
                h: 1920,
                w: 1080,
                ppi: 326,
                pxratio: 2.0,
                js: 1,
                language: "en".to_string(),
                carrier: Some("Verizon".to_string()),
                mccmnc: Some("310012".to_string()),
                connection_type: DeviceConnectionType::Wifi,
                geo: None,
            },
            ext: None,
            geo: Some(Geo {
                lat: Some(37.7749),
                lon: Some(-122.4194),
                accuracy: Some(10.6),
                country: Some("US".to_string()),
                city: Some("San Francisco".to_string()),
                zip: Some("94103".to_string()),
                utcoffset: Some(-8),
                lastfix: Some(1234567890),
            }),
            regs: None,
            segment: Some(sdk::Segment {
                id: None,
                uid: None,
                ext: None,
            }),
            session: Session {
                id: Uuid::new_v4(),
                launch_ts: 1234567890,
                launch_monotonic_ts: 1234567890,
                start_ts: 1234567890,
                start_monotonic_ts: 1234567890,
                ts: 1234567890,
                monotonic_ts: 1234567890,
                memory_warnings_ts: vec![],
                memory_warnings_monotonic_ts: vec![],
                ram_used: 1024,
                ram_size: 2048,
                storage_free: Some(512),
                storage_used: Some(256),
                battery: 80.5,
                cpu_usage: 10.6,
            },
            test: Some(false),
            tmax: Some(500),
            token: None,
            user: sdk::User {
                idfa: Some(Uuid::new_v4()),
                tracking_authorization_status: "authorized".to_string(),
                idfv: Some(Uuid::new_v4()),
                idg: Some(Uuid::new_v4()),
                coppa: None,
                consent: Some(HashMap::from([
                    ("meta".to_string(), json!({"consent": true})),
                    ("gdpr".to_string(), json!({"status": "granted"})),
                ])),
            },
        }
    }

    #[test]
    fn test_convert_session() {
        let id = Uuid::new_v4();

        let api_session = Session {
            id: id,
            launch_ts: 1234567890,
            launch_monotonic_ts: 1234567890,
            start_ts: 1234567890,
            start_monotonic_ts: 1234567890,
            ts: 1234567890,
            monotonic_ts: 1234567890,
            memory_warnings_ts: vec![1234567890],
            memory_warnings_monotonic_ts: vec![1234567890],
            ram_used: 1024,
            ram_size: 2048,
            storage_free: Some(512),
            storage_used: Some(256),
            battery: 80.5,
            cpu_usage: 10.6,
        };

        let bidon_session = convert_session(&api_session);

        assert_eq!(bidon_session.id, Some(id.to_string()));
        assert_eq!(bidon_session.launch_ts, Some(1234567890));
        assert_eq!(bidon_session.ram_used, Some(1024));
        assert_eq!(bidon_session.launch_monotonic_ts, Some(1234567890));
        assert_eq!(bidon_session.start_ts, Some(1234567890));
        assert_eq!(bidon_session.start_monotonic_ts, Some(1234567890));
        assert_eq!(bidon_session.ts, Some(1234567890));
        assert_eq!(bidon_session.monotonic_ts, Some(1234567890));
        assert_eq!(bidon_session.memory_warnings_ts, vec![1234567890]);
        assert_eq!(bidon_session.memory_warnings_monotonic_ts, vec![1234567890]);
        assert_eq!(bidon_session.ram_size, Some(2048));
        assert_eq!(bidon_session.storage_free, Some(512));
        assert_eq!(bidon_session.storage_used, Some(256));
        assert_eq!(bidon_session.battery, Some(80.5));
        assert_eq!(bidon_session.cpu_usage, Some(10.6));
    }

    #[test]
    fn test_convert_user() {
        let idfa = Uuid::new_v4();
        let idfv = Uuid::new_v4();
        let idg = Uuid::new_v4();

        let consent_map = HashMap::from([("meta".to_string(), json!({"consent": true}))]);

        let api_user = sdk::User {
            idfa: Some(idfa),
            tracking_authorization_status: "authorized".to_string(),
            idfv: Some(idfv),
            idg: Some(idg),
            consent: Some(consent_map.clone()),
            coppa: None,
        };

        let segment = Some(Segment {
            id: Some("segment_id".to_string()),
            uid: Some("segment_uid".to_string()),
            ext: None,
        });

        let user = convert_user(&api_user, segment.as_ref()).unwrap();

        assert_eq!(user.id, Some(idg.to_string()));
        assert_eq!(
            user.consent,
            Some(serde_json::to_string(&consent_map).unwrap())
        );

        let bidon_user = user.extension_set.extension_data(USER_EXT).unwrap();
        assert_eq!(bidon_user.idfa, Some(idfa.to_string()));
        assert_eq!(bidon_user.idfv, Some(idfv.to_string()));
        assert_eq!(
            bidon_user.tracking_authorization_status,
            Some("authorized".to_string())
        );
        assert_eq!(bidon_user.segments[0].id, Some("segment_id".to_string()));
        assert_eq!(bidon_user.segments[0].uid, Some("segment_uid".to_string()));
    }

    #[test]
    fn test_convert_app() {
        let api_app = App {
            version: "1.0".to_string(),
            bundle: "com.example.app".to_string(),
            key: "app_key".to_string(),
            framework_version: Some("1.0".to_string()),
            plugin_version: Some("1.0".to_string()),
            sdk_version: Some("1.0".to_string()),
            skadn: Some(vec!["skadn1".to_string(), "skadn2".to_string()]),
            framework: "".to_string(),
        };

        let app = convert_app(&api_app, &"1.0".to_string()).unwrap();

        assert_eq!(app.ver, Some("1.0".to_string()));
        assert_eq!(app.bundle, Some("com.example.app".to_string()));
        let bidon_app = app.extension_set.extension_data(APP_EXT).unwrap();
        assert_eq!(bidon_app.key, Some("app_key".to_string()));
        assert_eq!(bidon_app.framework, Some("1.0".to_string()));
        assert_eq!(bidon_app.framework_version, Some("1.0".to_string()));
        assert_eq!(bidon_app.plugin_version, Some("1.0".to_string()));
        assert_eq!(bidon_app.sdk_version, Some("1.0".to_string()));
        assert_eq!(
            bidon_app.skadn,
            vec!["skadn1".to_string(), "skadn2".to_string()]
        );
    }

    #[test]
    fn test_convert_regs() {
        let api_regs = sdk::Regulations {
            coppa: Some(true),
            gdpr: Some(true),
            us_privacy: Some("1YNN".to_string()),
            eu_privacy: Some("1".to_string()),
            iab: Some(HashMap::from([("key".to_string(), json!("value"))])),
        };

        let regs = convert_regs(&api_regs).unwrap();

        assert_eq!(regs.coppa, Some(true));
        assert_eq!(regs.gdpr, Some(true));
        let bidon_regs = regs.extension_set.extension_data(REGS_EXT).unwrap();
        assert_eq!(bidon_regs.us_privacy, Some("1YNN".to_string()));
        assert_eq!(bidon_regs.eu_privacy, Some("1".to_string()));
        assert_eq!(bidon_regs.iab, Some("{\"key\":\"value\"}".to_string()));
    }

    #[test]
    fn test_convert_ad_object_to_item() {
        let ad_object = sdk::AdObject {
            auction_id: Some("auction_id".to_string()),
            auction_key: Some("auction_key".to_string()),
            auction_configuration_id: Some(123i64),
            auction_configuration_uid: Some("auction_configuration_uid".to_string()),
            orientation: Some(AdObjectOrientation::Portrait),
            demands: HashMap::from([(
                "demand_key".to_string(),
                json!({
                    "token": "token_value",
                    "status": "status_value",
                    "token_finish_ts": 1234567890,
                    "token_start_ts": 1234567990
                }),
            )]),
            banner: Some(sdk::BannerAdObject {
                format: AdFormat::Banner,
            }),
            interstitial: Some(json!({"interstitial": "value"})),
            rewarded: Some(json!("rewarded".to_string())),
            auction_pricefloor: 1.0,
        };

        let item = convert_ad_object_to_item(&ad_object).unwrap();

        assert_eq!(item.id, Some("auction_id".to_string()));
        assert_eq!(item.flr, Some(1.0));
        assert_eq!(item.flrcur, Some("USD".to_string()));
        let placement = item.spec.unwrap();
        fn registry(extension: &'static dyn Extension) -> ExtensionRegistry {
            let mut registry = ExtensionRegistry::new();
            registry.register(extension);
            registry
        }
        let placement =
            crate::com::iabtechlab::adcom::v1::placement::Placement::decode_with_extensions(
                placement.as_slice(),
                &registry(PLACEMENT_EXT),
            )
            .unwrap();
        let bidon_ad_object = placement
            .extension_set
            .extension_data(PLACEMENT_EXT)
            .unwrap();
        assert_eq!(bidon_ad_object.auction_id, Some("auction_id".to_string()));
        assert_eq!(bidon_ad_object.auction_key, Some("auction_key".to_string()));
        assert_eq!(bidon_ad_object.auction_configuration_id, Some(123));
        assert_eq!(
            bidon_ad_object.auction_configuration_uid,
            Some("auction_configuration_uid".to_string())
        );
        assert_eq!(
            bidon_ad_object.orientation,
            Some(Orientation::Portrait as i32)
        );
        assert_eq!(
            bidon_ad_object.demands.get("demand_key").unwrap().token,
            Some("token_value".to_string())
        );
        assert_eq!(
            bidon_ad_object.demands.get("demand_key").unwrap().status,
            Some("status_value".to_string())
        );
        assert_eq!(
            bidon_ad_object
                .demands
                .get("demand_key")
                .unwrap()
                .token_finish_ts,
            Some(1234567890)
        );
        assert_eq!(
            bidon_ad_object
                .demands
                .get("demand_key")
                .unwrap()
                .token_start_ts,
            Some(1234567990)
        );
        assert_eq!(
            bidon_ad_object.banner.as_ref().unwrap().format,
            Some(crate::org::bidon::proto::v1::mediation::AdFormat::Banner as i32)
        );
        assert_eq!(
            bidon_ad_object.interstitial,
            Some("{\"interstitial\":\"value\"}".to_string())
        );
        assert_eq!(bidon_ad_object.rewarded, Some("\"rewarded\"".to_string()));
    }

    #[test]
    fn test_openrtb_to_auction_response() {
        // Create a mock OpenRTB response
        let response = openrtb::Response {
            id: Some("auction123".to_string()),
            bidid: None,
            nbr: None,
            seatbid: vec![openrtb::SeatBid {
                bid: vec![openrtb::Bid {
                    id: Some("bid1".to_string()),
                    item: Some("item1".to_string()),
                    price: Some(2.5),
                    cid: Some("demand1".to_string()),
                    extension_set: {
                        let mut ext = prost::ExtensionSet::default();
                        let bid_ext = mediation::BidExt {
                            label: Some("key123".to_string()),
                            bid_type: Some("bid_type".to_string()),
                            ext: HashMap::new(),
                        };
                        ext.set_extension_data(BID_EXT, bid_ext).unwrap();
                        ext
                    },
                    ..Default::default()
                }],
                seat: None,
                ..Default::default()
            }],
            extension_set: {
                let mut ext = prost::ExtensionSet::default();
                let auction_response_ext = mediation::AuctionResponseExt {
                    auction_id: Some("auction123".to_string()),
                    auction_configuration_id: Some(456),
                    auction_configuration_uid: Some("config789".to_string()),
                    token: Some("key123".to_string()),
                    auction_pricefloor: Some(1.0),
                    auction_timeout: Some(500),
                    external_win_notifications: Some(true),
                    segment: Some(mediation::Segment {
                        id: Some("segment_id".to_string()),
                        uid: Some("segment_uid".to_string()),
                        ext: None,
                    }),
                };
                ext.set_extension_data(AUCTION_RESPONSE_EXT, auction_response_ext)
                    .unwrap();
                ext
            },
            ..Default::default()
        };

        let openrtb = openrtb::Openrtb {
            ver: Some("3.0".to_string()),
            domainspec: Some("adcom".to_string()),
            domainver: Some("1.0".to_string()),
            payload_oneof: Some(openrtb::openrtb::PayloadOneof::Response(response)),
        };

        let auction_response = try_into(openrtb).unwrap();

        // Test assertions
        assert_eq!(auction_response.auction_id, "auction123");
        assert_eq!(auction_response.auction_configuration_id, 456);
        assert_eq!(auction_response.auction_configuration_uid, "config789");
        assert_eq!(auction_response.auction_pricefloor, 1.0);
        assert_eq!(auction_response.auction_timeout, 500);
        assert_eq!(auction_response.token, "key123");

        // Test ad units
        assert!(!auction_response.ad_units.is_empty());
        let ad_unit = &auction_response.ad_units[0];
        assert_eq!(ad_unit.label, "key123");
        assert_eq!(ad_unit.uid, "item1");
        assert_eq!(ad_unit.demand_id, "demand1");
        assert_eq!(ad_unit.pricefloor, Some(2.5));
    }

    #[test]
    fn test_convert_device() {
        let request = create_test_auction_request();
        let device = convert_device(&request).unwrap();

        // Test standard fields
        assert_eq!(device.r#type, Some(adcom::enums::DeviceType::Phone as i32));
        assert_eq!(device.ua, Some("Mozilla/5.0".to_string()));
        assert_eq!(device.make, Some("Apple".to_string()));
        assert_eq!(device.model, Some("iPhone".to_string()));
        assert_eq!(device.os, Some(OperatingSystem::Ios as i32));
        assert_eq!(device.osv, Some("14.4".to_string()));
        assert_eq!(device.hwv, Some("A14".to_string()));
        assert_eq!(device.h, Some(1920));
        assert_eq!(device.w, Some(1080));
        assert_eq!(device.ppi, Some(326));
        assert_eq!(device.pxratio, Some(2.0));
        assert_eq!(device.js, Some(true));
        assert_eq!(device.lang, Some("en".to_string()));
        assert_eq!(device.carrier, Some("Verizon".to_string()));
        assert_eq!(device.mccmnc, Some("310012".to_string()));
        assert_eq!(
            device.contype,
            Some(adcom::enums::ConnectionType::Wifi as i32)
        );

        // Test geo fields
        let geo = device.geo.unwrap();
        assert_eq!(geo.lat, Some(37.7749));
        assert_eq!(geo.lon, Some(-122.4194));
        assert_eq!(geo.accur, Some(10)); // Converted from f64 to i32
        assert_eq!(geo.country, Some("US".to_string()));
        assert_eq!(geo.city, Some("San Francisco".to_string()));
        assert_eq!(geo.zip, Some("94103".to_string()));
        assert_eq!(geo.utcoffset, Some(-8));
        assert_eq!(geo.lastfix, Some(1234567890));

        // Test device extension fields
        let device_ext = device.extension_set.extension_data(DEVICE_EXT).unwrap();
        assert!(device_ext.id.is_some());
        assert_eq!(device_ext.launch_ts, Some(1234567890));
        assert_eq!(device_ext.launch_monotonic_ts, Some(1234567890));
        assert_eq!(device_ext.start_ts, Some(1234567890));
        assert_eq!(device_ext.start_monotonic_ts, Some(1234567890));
        assert_eq!(device_ext.ts, Some(1234567890));
        assert_eq!(device_ext.monotonic_ts, Some(1234567890));
        assert!(device_ext.memory_warnings_ts.is_empty());
        assert!(device_ext.memory_warnings_monotonic_ts.is_empty());
        assert_eq!(device_ext.ram_used, Some(1024));
        assert_eq!(device_ext.ram_size, Some(2048));
        assert_eq!(device_ext.storage_free, Some(512));
        assert_eq!(device_ext.storage_used, Some(256));
        assert_eq!(device_ext.battery, Some(80.5));
        assert_eq!(device_ext.cpu_usage, Some(10.6));
    }
}
