use crate::com::iabtechlab::adcom::v1 as adcom;
use crate::com::iabtechlab::adcom::v1::context::DistributionChannel;
use crate::com::iabtechlab::adcom::v1::enums::{ConnectionType, OperatingSystem};
use crate::com::iabtechlab::adcom::v1::placement::placement::DisplayPlacement;
use crate::com::iabtechlab::openrtb::v3 as openrtb;
use crate::com::iabtechlab::openrtb::v3::AuctionType;
use crate::org::bidon::proto::v1::context::Context;
use crate::org::bidon::proto::v1::mediation::{self, RequestExt, SdkAdapter};
use crate::sdk;
use anyhow::{anyhow, Result};
use prost::{Extendable, Message};
use serde_json::Value;
use std::collections::HashMap;
use std::net::IpAddr;

pub(crate) fn try_from(
    request: &sdk::AuctionRequest,
    bidon_version: String,
    ip: IpAddr,
    ad_type: sdk::GetAuctionAdTypeParameter,
) -> Result<openrtb::Openrtb> {
    let context = convert_context(&request, bidon_version, ip)?;

    // Convert AuctionRequest to Openrtb::Request
    let mut openrtb_request = openrtb::Request {
        id: request.ad_object.auction_id.to_owned(),
        test: request.test,
        tmax: request.tmax.map(|t| t as u32),
        at: Some(AuctionType::FirstPrice as i32),
        context: context.encode_to_vec().into(),
        item: vec![convert_ad_object_to_item(&request.ad_object, ad_type)?],
        ..Default::default()
    };

    let request_ext = RequestExt {
        ad_type: Some(convert_ad_type(&ad_type) as i32),
        adapters: request
            .adapters
            .clone()
            .into_iter()
            .map(|(k, v)| (k, convert_adapter(&v)))
            .collect(),
        ext: request.ext.clone(),
    };

    openrtb_request.set_extension_data(mediation::REQUEST_EXT, request_ext)?;

    // Create Openrtb instance with the converted request
    Ok(openrtb::Openrtb {
        ver: Some("3.0".to_string()),
        domainspec: Some("domain_spec".to_string()),
        domainver: Some("domain_version".to_string()),
        payload_oneof: Some(openrtb::openrtb::PayloadOneof::Request(openrtb_request)),
    })
}

fn convert_context(
    request: &sdk::AuctionRequest,
    bidon_version: String,
    ip: IpAddr,
) -> Result<Context> {
    // Create the AdCOM Context message
    Ok(Context {
        distribution_channel: DistributionChannel {
            channel_oneof: Some(adcom::context::distribution_channel::ChannelOneof::App(
                convert_app(&request.app, bidon_version)?,
            )),
            ..Default::default()
        }
        .into(),
        device: convert_device(&request.device, &request.session, request.geo.as_ref(), ip)?.into(),
        user: convert_user(&request.user, request.segment.as_ref())?.into(),
        regs: match request.regs.as_ref() {
            Some(t) => convert_regs(t)?.into(),
            None => None,
        },
        restrictions: None, // TODO
    })
}

fn convert_app(
    app: &sdk::App,
    bidon_version: String,
) -> Result<adcom::context::distribution_channel::App> {
    let mut adcom_app = adcom::context::distribution_channel::App {
        ver: app.version.clone().into(),
        keywords: None,
        paid: None,
        bundle: app.bundle.clone().into(),
        ..Default::default()
    };

    let app_ext = mediation::AppExt {
        key: app.key.clone().into(),
        framework: app.framework_version.clone(),
        framework_version: app.framework_version.clone(),
        plugin_version: app.plugin_version.clone(),
        sdk_version: app.sdk_version.clone(),
        skadn: app.skadn.clone().unwrap_or_default(),
        bidon_version: Some(bidon_version.clone()),
    };

    adcom_app.set_extension_data(mediation::APP_EXT, app_ext)?;
    Ok(adcom_app)
}

fn convert_device(
    device: &sdk::Device,
    session: &sdk::Session,
    geo: Option<&sdk::Geo>,
    ip: IpAddr,
) -> Result<adcom::context::Device> {
    let mut adcom_device = adcom::context::Device {
        // Map standard fields
        r#type: convert_device_type(device.device_type).map(Into::into),
        ua: device.ua.clone().into(),
        make: device.make.clone().into(),
        model: device.model.clone().into(),
        os: Some(Into::into(convert_os(device.os.clone()))),
        osv: device.osv.clone().into(),
        hwv: device.hwv.clone().into(),
        h: device.h.into(),
        w: device.w.into(),
        ppi: device.ppi.into(),
        pxratio: (device.pxratio as f32).into(), // TODO validate conversion
        js: Some(device.js != 0),
        lang: device.language.clone().into(),
        carrier: device.clone().carrier,
        mccmnc: device.clone().mccmnc,
        contype: Some(Into::into(convert_connection_type(device.connection_type))),
        geo: geo.map(convert_geo),
        ip: Some(ip.to_string()),
        ipv6: match ip {
            IpAddr::V4(_) => None,
            IpAddr::V6(ip) => Some(ip.to_string()),
        },
        ..Default::default()
    };

    adcom_device.set_extension_data(mediation::DEVICE_EXT, convert_session(session))?;

    Ok(adcom_device)
}

fn convert_device_type(device_type: Option<sdk::DeviceType>) -> Option<adcom::enums::DeviceType> {
    device_type
        .map(|dt| match dt {
            sdk::DeviceType::Phone => adcom::enums::DeviceType::Phone,
            sdk::DeviceType::Tablet => adcom::enums::DeviceType::Tablet,
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

fn convert_connection_type(conn_type: sdk::DeviceConnectionType) -> ConnectionType {
    match conn_type {
        sdk::DeviceConnectionType::Ethernet => ConnectionType::Wired,
        sdk::DeviceConnectionType::Wifi => ConnectionType::Wifi,
        sdk::DeviceConnectionType::CellularUnknown => ConnectionType::CellUnknown,
        sdk::DeviceConnectionType::Cellular => ConnectionType::CellUnknown,
        sdk::DeviceConnectionType::Cellular2G => ConnectionType::Cell2g,
        sdk::DeviceConnectionType::Cellular3G => ConnectionType::Cell3g,
        sdk::DeviceConnectionType::Cellular4G => ConnectionType::Cell4g,
        sdk::DeviceConnectionType::Cellular5G => ConnectionType::Cell5g,
    }
}

fn convert_geo(geo: &sdk::Geo) -> adcom::context::Geo {
    adcom::context::Geo {
        r#type: Some(adcom::enums::LocationType::Unknown as i32), // TODO
        lat: geo.lat.map(|t| t as f32),
        lon: geo.lon.map(|t| t as f32),
        accur: geo.accuracy.map(|t| (t as i32)), // TODO check accuracy conversion. We convert it from f64 to i32 here.
        country: geo.country.clone(),
        city: geo.city.clone(),
        zip: geo.zip.clone(),
        utcoffset: geo.utcoffset,
        lastfix: geo.lastfix,
        ..Default::default()
    }
}

fn convert_user(user: &sdk::User, segment: Option<&sdk::Segment>) -> Result<adcom::context::User> {
    let mut adcom_user = adcom::context::User {
        id: user.idg.map(|uuid| uuid.to_string()),
        consent: serde_json::to_string(&user.consent).ok(),
        ..Default::default()
    };

    let user_ext = mediation::UserExt {
        idfa: user.idfa.map(|uuid| uuid.to_string()),
        tracking_authorization_status: Some(user.tracking_authorization_status.clone()),
        idfv: user.idfv.map(|uuid| uuid.to_string()),
        idg: user.idg.map(|uuid| uuid.to_string()),
        segments: segment.into_iter().map(convert_segment).collect(),
    };

    adcom_user.set_extension_data(mediation::USER_EXT, user_ext)?;

    Ok(adcom_user)
}

fn convert_adapter(adapter: &sdk::Adapter) -> SdkAdapter {
    SdkAdapter {
        version: Some(adapter.version.clone()),
        sdk_version: Some(adapter.sdk_version.clone()),
    }
}

fn convert_ad_type(ad_type: &sdk::GetAuctionAdTypeParameter) -> mediation::AdType {
    match ad_type {
        sdk::GetAuctionAdTypeParameter::Banner => mediation::AdType::Banner,
        sdk::GetAuctionAdTypeParameter::Interstitial => mediation::AdType::Interstitial,
        sdk::GetAuctionAdTypeParameter::Rewarded => mediation::AdType::Rewarded,
    }
}

fn convert_segment(segment: &sdk::Segment) -> mediation::Segment {
    mediation::Segment {
        id: segment.id.clone(),
        uid: segment.uid.clone(),
        ext: segment.ext.clone(),
    }
}

fn convert_session(session: &sdk::Session) -> mediation::DeviceExt {
    mediation::DeviceExt {
        id: Some(session.id.to_string().clone()),
        launch_ts: Some(session.launch_ts),
        launch_monotonic_ts: session.launch_monotonic_ts.into(),
        start_ts: Some(session.start_ts),
        start_monotonic_ts: Some(session.start_monotonic_ts),
        ts: Some(session.ts),
        monotonic_ts: Some(session.monotonic_ts),
        memory_warnings_ts: session.memory_warnings_ts.clone(),
        memory_warnings_monotonic_ts: session.memory_warnings_monotonic_ts.clone(),
        ram_used: Some(session.ram_used),
        ram_size: Some(session.ram_size),
        storage_free: session.storage_free,
        storage_used: session.storage_used,
        battery: Some(session.battery),
        cpu_usage: Some(session.cpu_usage),
    }
}

fn convert_regs(regs: &sdk::Regulations) -> Result<adcom::context::Regs> {
    let mut adcom_regs = adcom::context::Regs {
        coppa: regs.coppa,
        gdpr: regs.gdpr,
        ..Default::default()
    };

    let regs_ext = mediation::RegsExt {
        us_privacy: regs.us_privacy.clone(),
        eu_privacy: regs.eu_privacy.clone(),
        iab: match regs.iab.as_ref() {
            Some(t) => Some(convert_iab(t)?),
            None => None,
        },
    };

    adcom_regs.set_extension_data(mediation::REGS_EXT, regs_ext)?;

    Ok(adcom_regs)
}

fn convert_iab(iab_json: &HashMap<String, Value>) -> Result<String> {
    serde_json::to_string(&iab_json).map_err(Into::into)
}

fn convert_ad_object_to_item(
    ad_object: &sdk::AdObject,
    ad_type: sdk::GetAuctionAdTypeParameter,
) -> Result<openrtb::Item> {
    let mut placement = adcom::placement::Placement {
        display: match (ad_object.orientation, ad_object.banner.as_ref()) {
            (None, None) => None,
            (orientation, banner) => Some(DisplayPlacement {
                extension_set: {
                    let mut ext = prost::ExtensionSet::default();
                    let dpe = mediation::DisplayPlacementExt {
                        orientation: orientation.map(|o| convert_ad_orientation(&o) as i32),
                        format: banner.map(|b| convert_banner_format(b.format) as i32),
                    };
                    ext.set_extension_data(mediation::DISPLAY_PLACEMENT_EXT, dpe)?;
                    ext
                },
                ..Default::default()
            }),
        },
        ..Default::default()
    };

    // Convert based on ad type and orientation
    match ad_type {
        sdk::GetAuctionAdTypeParameter::Banner => {
            // Regular banner. 0 is false, 1 is true
            placement.display.get_or_insert(Default::default()).instl = Some(0);
        }
        sdk::GetAuctionAdTypeParameter::Interstitial => {
            // Interstitial can be either display or video. 0 is false, 1 is true
            placement.display.get_or_insert(Default::default()).instl = Some(1);
            placement.video = Some(adcom::placement::placement::VideoPlacement {
                ptype: Some(adcom::enums::VideoPlacementSubtype::Interstitial as i32),
                ..Default::default()
            });
        }
        sdk::GetAuctionAdTypeParameter::Rewarded => {
            // Rewarded can be either display or video.
            placement.reward = Some(true);
            placement.video = Some(adcom::placement::placement::VideoPlacement {
                ..Default::default()
            });
        }
    }

    // Add extension data for PlacementExt
    let placement_ext = mediation::PlacementExt {
        auction_id: ad_object.auction_id.clone(),
        auction_key: ad_object.auction_key.clone(),
        auction_configuration_id: ad_object.auction_configuration_id,
        auction_configuration_uid: ad_object.auction_configuration_uid.clone(),
        demands: convert_demand(&ad_object.demands)?,
        ..Default::default()
    };

    placement.set_extension_data(mediation::PLACEMENT_EXT, placement_ext)?;

    Ok(openrtb::Item {
        id: ad_object.auction_id.clone(),
        flr: Some(ad_object.auction_pricefloor as f32),
        flrcur: Some("USD".to_string()),
        spec: placement.encode_to_vec().into(),
        ..Default::default()
    })
}

fn convert_banner_format(format: sdk::AdFormat) -> mediation::AdFormat {
    match format {
        sdk::AdFormat::Banner => mediation::AdFormat::Banner,
        sdk::AdFormat::Leaderboard => mediation::AdFormat::Leaderboard,
        sdk::AdFormat::Mrec => mediation::AdFormat::Mrec,
        sdk::AdFormat::Adaptive => mediation::AdFormat::Adaptive,
    }
}

fn convert_ad_orientation(orientation: &sdk::AdObjectOrientation) -> mediation::Orientation {
    match orientation {
        sdk::AdObjectOrientation::Portrait => mediation::Orientation::Portrait,
        sdk::AdObjectOrientation::Landscape => mediation::Orientation::Landscape,
    }
}

fn convert_demand(demand: &HashMap<String, Value>) -> Result<HashMap<String, mediation::Demand>> {
    let mut demands = HashMap::new();

    // TODO: mb we should preserve JSON structure?
    for (key, value) in demand {
        let map = value
            .as_object()
            .ok_or(anyhow!("Demand value is not an object: {}", value))?;
        let mediation_demand = mediation::Demand {
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
        demands.insert(key.clone(), mediation_demand);
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

    // Extract auction configuration from response extensions
    let auction_ext = response
        .extension_set
        .extension_data(mediation::AUCTION_RESPONSE_EXT)
        .map_err(|_| anyhow!("Missing mediation ad object extension in response"))?;

    // Extract bid information from the response
    let mut ad_units = Vec::new();
    let mut no_bids = Vec::new();

    for seatbid in response.seatbid {
        for bid in seatbid.bid {
            // Extract bid extension data
            let bid_ext = bid
                .extension_set
                .extension_data(mediation::BID_EXT)
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
            if bid.price.unwrap_or_default() as f64
                > auction_ext.auction_pricefloor.unwrap_or_default()
            {
                ad_units.push(ad_unit);
            } else {
                no_bids.push(ad_unit);
            }
        }
    }

    let auction_response = sdk::AuctionResponse {
        ad_units,
        auction_id: response.id.unwrap_or_default(),
        no_bids: Some(no_bids),
        token: auction_ext.token.clone().unwrap_or_default(),
        external_win_notifications: auction_ext.external_win_notifications.unwrap_or_default(),
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
    use prost::ExtensionRegistry;
    use serde_json::json;
    use std::collections::HashMap;
    use std::io::Cursor;
    use std::net::Ipv4Addr;
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
            app: sdk::App {
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
                device_type: Some(sdk::DeviceType::Phone),
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
                connection_type: sdk::DeviceConnectionType::Wifi,
                geo: None,
            },
            ext: None,
            geo: Some(sdk::Geo {
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
            session: sdk::Session {
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
    fn test_convert_device() {
        let request = create_test_auction_request();
        let ip = IpAddr::V4(Ipv4Addr::new(127, 0, 0, 1));
        let adcom_device =
            convert_device(&request.device, &request.session, request.geo.as_ref(), ip).unwrap();

        // Test standard fields
        assert_eq!(
            adcom_device.r#type,
            Some(adcom::enums::DeviceType::Phone as i32)
        );
        assert_eq!(adcom_device.ua, Some("Mozilla/5.0".to_string()));
        assert_eq!(adcom_device.make, Some("Apple".to_string()));
        assert_eq!(adcom_device.model, Some("iPhone".to_string()));
        assert_eq!(adcom_device.os, Some(OperatingSystem::Ios as i32));
        assert_eq!(adcom_device.osv, Some("14.4".to_string()));
        assert_eq!(adcom_device.hwv, Some("A14".to_string()));
        assert_eq!(adcom_device.h, Some(1920));
        assert_eq!(adcom_device.w, Some(1080));
        assert_eq!(adcom_device.ppi, Some(326));
        assert_eq!(adcom_device.pxratio, Some(2.0));
        assert_eq!(adcom_device.js, Some(true));
        assert_eq!(adcom_device.lang, Some("en".to_string()));
        assert_eq!(adcom_device.carrier, Some("Verizon".to_string()));
        assert_eq!(adcom_device.mccmnc, Some("310012".to_string()));
        assert_eq!(
            adcom_device.contype,
            Some(adcom::enums::ConnectionType::Wifi as i32)
        );

        // Test geo fields
        let geo = adcom_device.geo.unwrap();
        assert_eq!(geo.lat, Some(37.7749));
        assert_eq!(geo.lon, Some(-122.4194));
        assert_eq!(geo.accur, Some(10)); // Converted from f64 to i32
        assert_eq!(geo.country, Some("US".to_string()));
        assert_eq!(geo.city, Some("San Francisco".to_string()));
        assert_eq!(geo.zip, Some("94103".to_string()));
        assert_eq!(geo.utcoffset, Some(-8));
        assert_eq!(geo.lastfix, Some(1234567890));

        // Test device extension fields
        let device_ext = adcom_device
            .extension_set
            .extension_data(mediation::DEVICE_EXT)
            .unwrap();
        assert_eq!(device_ext.id, Some(request.session.id.to_string()));
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

    #[test]
    fn test_convert_user() {
        let idfa = Uuid::new_v4();
        let idfv = Uuid::new_v4();
        let idg = Uuid::new_v4();

        let user = sdk::User {
            idfa: Some(idfa),
            tracking_authorization_status: "authorized".to_string(),
            idfv: Some(idfv),
            idg: Some(idg),
            consent: Some(HashMap::from([(
                "meta".to_string(),
                json!({"consent": true}),
            )])),
            coppa: None,
        };

        let segment = Some(sdk::Segment {
            id: Some("segment_id".to_string()),
            uid: Some("segment_uid".to_string()),
            ext: None,
        });

        let adcom_user = convert_user(&user, segment.as_ref()).unwrap();

        assert_eq!(adcom_user.id, Some(idg.to_string()));
        let user_ext = adcom_user
            .extension_set
            .extension_data(mediation::USER_EXT)
            .unwrap();
        assert_eq!(user_ext.idfa, Some(idfa.to_string()));
        assert_eq!(user_ext.idfv, Some(idfv.to_string()));
        assert_eq!(
            user_ext.tracking_authorization_status,
            Some("authorized".to_string())
        );
        assert_eq!(
            adcom_user.consent,
            Some("{\"meta\":{\"consent\":true}}".to_string())
        );
        assert_eq!(user_ext.segments[0].id, Some("segment_id".to_string()));
        assert_eq!(user_ext.segments[0].uid, Some("segment_uid".to_string()));
    }

    #[test]
    fn test_convert_app() {
        let app = sdk::App {
            version: "1.0".to_string(),
            bundle: "com.example.app".to_string(),
            key: "app_key".to_string(),
            framework_version: Some("1.0".to_string()),
            plugin_version: Some("1.0".to_string()),
            sdk_version: Some("1.0".to_string()),
            skadn: Some(vec!["skadn1".to_string(), "skadn2".to_string()]),
            framework: "".to_string(),
        };

        let adcom_app = convert_app(&app, "1.0".to_string()).unwrap();

        assert_eq!(adcom_app.ver, Some("1.0".to_string()));
        assert_eq!(adcom_app.bundle, Some("com.example.app".to_string()));

        let app_ext = adcom_app
            .extension_set
            .extension_data(mediation::APP_EXT)
            .unwrap();
        assert_eq!(app_ext.key, Some("app_key".to_string()));
        assert_eq!(app_ext.framework, Some("1.0".to_string()));
        assert_eq!(app_ext.framework_version, Some("1.0".to_string()));
        assert_eq!(app_ext.plugin_version, Some("1.0".to_string()));
        assert_eq!(app_ext.sdk_version, Some("1.0".to_string()));
        assert_eq!(
            app_ext.skadn,
            vec!["skadn1".to_string(), "skadn2".to_string()]
        );
    }

    #[test]
    fn test_convert_regs() {
        let regs = sdk::Regulations {
            coppa: Some(true),
            gdpr: Some(true),
            us_privacy: Some("1YNN".to_string()),
            eu_privacy: Some("1".to_string()),
            iab: Some(HashMap::from([("key".to_string(), json!("value"))])),
        };

        let adcom_regs = convert_regs(&regs).unwrap();

        assert_eq!(adcom_regs.coppa, Some(true));
        assert_eq!(adcom_regs.gdpr, Some(true));
        let regs_ext = adcom_regs
            .extension_set
            .extension_data(mediation::REGS_EXT)
            .unwrap();
        assert_eq!(regs_ext.us_privacy, Some("1YNN".to_string()));
        assert_eq!(regs_ext.eu_privacy, Some("1".to_string()));
        assert_eq!(regs_ext.iab, Some("{\"key\":\"value\"}".to_string()));
    }

    #[test]
    fn test_convert_ad_object_to_item() {
        // Test banner ad object
        let banner_ad_object = sdk::AdObject {
            auction_id: Some("auction_id".to_string()),
            auction_key: Some("auction_key".to_string()),
            auction_configuration_id: Some(123i64),
            auction_configuration_uid: Some("auction_configuration_uid".to_string()),
            orientation: Some(sdk::AdObjectOrientation::Portrait),
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
                format: sdk::AdFormat::Banner,
            }),
            interstitial: None,
            rewarded: None,
            auction_pricefloor: 1.0,
        };

        let item =
            convert_ad_object_to_item(&banner_ad_object, sdk::GetAuctionAdTypeParameter::Banner)
                .unwrap();

        let mut registry = ExtensionRegistry::new();
        registry.register(mediation::PLACEMENT_EXT);
        registry.register(mediation::DISPLAY_PLACEMENT_EXT);

        let placement = adcom::placement::Placement::decode_with_extensions(
            &mut Cursor::new(item.spec.unwrap()),
            &registry,
        )
        .unwrap();

        let placement_ext = placement
            .extension_set
            .extension_data(mediation::PLACEMENT_EXT)
            .unwrap();

        // Test basic item fields
        assert_eq!(item.id, Some("auction_id".to_string()));
        assert_eq!(item.flr, Some(1.0));
        assert_eq!(item.flrcur, Some("USD".to_string()));

        // Test placement extension fields
        assert_eq!(placement_ext.auction_id, Some("auction_id".to_string()));
        assert_eq!(placement_ext.auction_key, Some("auction_key".to_string()));
        assert_eq!(placement_ext.auction_configuration_id, Some(123));
        assert_eq!(
            placement_ext.auction_configuration_uid,
            Some("auction_configuration_uid".to_string())
        );

        // Test demands
        let demand = placement_ext.demands.get("demand_key").unwrap();
        assert_eq!(demand.token, Some("token_value".to_string()));
        assert_eq!(demand.status, Some("status_value".to_string()));
        assert_eq!(demand.token_finish_ts, Some(1234567890));
        assert_eq!(demand.token_start_ts, Some(1234567990));

        // Test display placement fields
        let display = placement.display.unwrap();
        let display_ext = display
            .extension_set
            .extension_data(mediation::DISPLAY_PLACEMENT_EXT)
            .unwrap();

        assert_eq!(display.instl, Some(0)); // Regular banner
        assert_eq!(
            display_ext.orientation,
            Some(mediation::Orientation::Portrait as i32)
        );
        assert_eq!(display_ext.format, Some(mediation::AdFormat::Banner as i32));
        assert_eq!(display.instl, Some(0));
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
                        ext.set_extension_data(mediation::BID_EXT, bid_ext).unwrap();
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
                ext.set_extension_data(mediation::AUCTION_RESPONSE_EXT, auction_response_ext)
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
    fn test_request_ext() {
        // Create a test auction request with adapters
        let mut request = create_test_auction_request();

        // Add test adapters
        request.adapters = HashMap::from([
            (
                "adapter1".to_string(),
                sdk::Adapter {
                    version: "1.0".to_string(),
                    sdk_version: "2.0".to_string(),
                },
            ),
            (
                "adapter2".to_string(),
                sdk::Adapter {
                    version: "3.0".to_string(),
                    sdk_version: "4.0".to_string(),
                },
            ),
        ]);

        // Add test ad type
        let ad_type = sdk::GetAuctionAdTypeParameter::Banner;

        // Convert to OpenRTB
        let openrtb = try_from(
            &request,
            "1.0".to_string(),
            IpAddr::V4(Ipv4Addr::new(127, 0, 0, 1)),
            ad_type,
        )
        .unwrap();

        // Extract the request from OpenRTB
        let openrtb_request = match openrtb.payload_oneof.unwrap() {
            openrtb::openrtb::PayloadOneof::Request(req) => req,
            _ => panic!("Expected Request payload"),
        };

        // Get the request extension
        let request_ext = openrtb_request
            .extension_set
            .extension_data(mediation::REQUEST_EXT)
            .unwrap();

        // Test ad type
        assert_eq!(request_ext.ad_type, Some(mediation::AdType::Banner as i32));

        // Test adapters
        assert_eq!(request_ext.adapters.len(), 2);

        let adapter1 = request_ext.adapters.get("adapter1").unwrap();
        assert_eq!(adapter1.version, Some("1.0".to_string()));
        assert_eq!(adapter1.sdk_version, Some("2.0".to_string()));

        let adapter2 = request_ext.adapters.get("adapter2").unwrap();
        assert_eq!(adapter2.version, Some("3.0".to_string()));
        assert_eq!(adapter2.sdk_version, Some("4.0".to_string()));

        // Test other ad types
        let interstitial_request = try_from(
            &request,
            "1.0".to_string(),
            IpAddr::V4(Ipv4Addr::new(127, 0, 0, 1)),
            sdk::GetAuctionAdTypeParameter::Interstitial,
        )
        .unwrap();

        let interstitial_request = match interstitial_request.payload_oneof.unwrap() {
            openrtb::openrtb::PayloadOneof::Request(req) => req,
            _ => panic!("Expected Request payload"),
        };

        let interstitial_ext = interstitial_request
            .extension_set
            .extension_data(mediation::REQUEST_EXT)
            .unwrap();

        assert_eq!(
            interstitial_ext.ad_type,
            Some(mediation::AdType::Interstitial as i32)
        );

        let rewarded_request = try_from(
            &request,
            "1.0".to_string(),
            IpAddr::V4(Ipv4Addr::new(127, 0, 0, 1)),
            sdk::GetAuctionAdTypeParameter::Rewarded,
        )
        .unwrap();

        let rewarded_request = match rewarded_request.payload_oneof.unwrap() {
            openrtb::openrtb::PayloadOneof::Request(req) => req,
            _ => panic!("Expected Request payload"),
        };

        let rewarded_ext = rewarded_request
            .extension_set
            .extension_data(mediation::REQUEST_EXT)
            .unwrap();

        assert_eq!(
            rewarded_ext.ad_type,
            Some(mediation::AdType::Rewarded as i32)
        );
    }
}
