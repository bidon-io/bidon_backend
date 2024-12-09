use crate::com::iabtechlab::adcom::v1 as adcom;
use crate::com::iabtechlab::adcom::v1::context::DistributionChannel;
use crate::com::iabtechlab::adcom::v1::enums::{ConnectionType, OperatingSystem};
use crate::com::iabtechlab::openrtb::v3 as openrtb;
use crate::com::iabtechlab::openrtb::v3::AuctionType;
use crate::org::bidon::proto::v1::context::Context;
use crate::org::bidon::proto::v1::mediation;
use crate::sdk;
use crate::sdk::AdObjectOrientation;
use anyhow::{anyhow, Result};
use prost::{Extendable, Message};
use serde_json::Value;
use std::collections::HashMap;

pub(crate) fn try_from(
    request: sdk::AuctionRequest,
    bidon_version: String,
) -> Result<openrtb::Openrtb> {
    let context = convert_context(&request, bidon_version)?;

    // Convert AuctionRequest to Openrtb::Request
    let openrtb_request = openrtb::Request {
        id: request.ad_object.auction_id.to_owned(),
        test: request.test,
        tmax: request.tmax.map(|t| t as u32),
        at: Some(AuctionType::FirstPrice as i32),
        context: context.encode_to_vec().into(),
        item: vec![convert_ad_object_to_item(&request.ad_object)?],
        ..Default::default()
    };

    // Create Openrtb instance with the converted request
    Ok(openrtb::Openrtb {
        ver: Some("3.0".to_string()),
        domainspec: Some("domain_spec".to_string()),
        domainver: Some("domain_version".to_string()),
        payload_oneof: Some(openrtb::openrtb::PayloadOneof::Request(openrtb_request)),
    })
}

fn convert_context(request: &sdk::AuctionRequest, bidon_version: String) -> Result<Context> {
    // Create the AdCOM Context message
    Ok(Context {
        distribution_channel: DistributionChannel {
            channel_oneof: Some(adcom::context::distribution_channel::ChannelOneof::App(
                convert_app(&request.app, bidon_version)?,
            )),
            ..Default::default()
        }
        .into(),
        device: convert_device(&request.device, &request.session, request.geo.as_ref())?.into(),
        user: convert_user(&request.user, request.segment.as_ref())?.into(),
        regs: match request.regs.as_ref() {
            Some(t) => convert_regs(&t)?.into(),
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
        skadn: app.skadn.clone().unwrap_or(vec![]),
        bidon_version: Some(bidon_version.clone()),
    };

    adcom_app.set_extension_data(mediation::APP_EXT, app_ext)?;
    Ok(adcom_app)
}

fn convert_device(
    device: &sdk::Device,
    session: &sdk::Session,
    geo: Option<&sdk::Geo>,
) -> Result<adcom::context::Device> {
    let mut adcom_device = adcom::context::Device {
        // Map standard fields
        r#type: convert_device_type(device.device_type.clone()).map(Into::into),
        ua: device.ua.clone().into(),
        make: device.make.clone().into(),
        model: device.model.clone().into(),
        os: Some(convert_os(device.os.clone())).map(Into::into),
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
        contype: Some(convert_connection_type(device.connection_type)).map(Into::into),
        geo: geo.map(convert_geo),
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
        country: geo.country.clone().into(),
        city: geo.city.clone().into(),
        zip: geo.zip.clone().into(),
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
        gdpr: regs.gdpr.clone(),
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

fn convert_ad_object_to_item(ad_object: &sdk::AdObject) -> Result<openrtb::Item> {
    // TODO: model AdObject with Placement
    let mut placement = adcom::placement::Placement {
        ..Default::default()
    };

    let placement_ext = mediation::PlacementExt {
        auction_id: ad_object.auction_id.clone(), // Unique Request ID, same as OpenRtb.Request.id
        auction_key: ad_object.auction_key.clone(), // Generated key for the auction request
        auction_configuration_id: ad_object.auction_configuration_id.clone(), // Deprecated: ID of the auction configuration
        auction_configuration_uid: ad_object.auction_configuration_uid.clone(), // UID of the auction configuration
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
        interstitial: ad_object.interstitial.as_ref().map(|i| i.to_string()), // TODO: remove String
        rewarded: ad_object.rewarded.as_ref().map(|r| r.to_string()),         // TODO: remove String
    };

    placement.set_extension_data(mediation::PLACEMENT_EXT, placement_ext)?;

    let item = openrtb::Item {
        id: ad_object.auction_id.clone(),
        flr: Some(ad_object.auction_pricefloor as f32),
        flrcur: Some("USD".to_string()), // TODO
        spec: placement.encode_to_vec().into(),
        ..Default::default()
    };

    Ok(item)
}

fn convert_ad_orientation(orientation: &AdObjectOrientation) -> mediation::Orientation {
    match orientation {
        AdObjectOrientation::Portrait => mediation::Orientation::Portrait,
        AdObjectOrientation::Landscape => mediation::Orientation::Landscape,
    }
}

fn convert_banner_ad(banner: &sdk::BannerAdObject) -> mediation::BannerAd {
    mediation::BannerAd {
        format: Some(convert_ad_format(banner.format) as i32),
    }
}

fn convert_ad_format(format: sdk::AdFormat) -> mediation::AdFormat {
    match format {
        sdk::AdFormat::Banner => mediation::AdFormat::Banner,
        sdk::AdFormat::Leaderboard => mediation::AdFormat::Leaderboard,
        sdk::AdFormat::Mrec => mediation::AdFormat::Mrec,
        sdk::AdFormat::Adaptive => mediation::AdFormat::Adaptive,
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

#[cfg(test)]
mod tests {
    use super::*;
    use prost::ExtensionRegistry;
    use serde_json::json;
    use std::collections::HashMap;
    use std::io::Cursor;
    use uuid::Uuid;

    #[test]
    fn test_convert_device() {
        let device = sdk::Device {
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
        };

        let id = Uuid::new_v4();

        let session = sdk::Session {
            id,
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

        let geo = Some(sdk::Geo {
            lat: Some(37.7749),
            lon: Some(-122.4194),
            accuracy: Some(10.6),
            country: Some("US".to_string()),
            city: Some("San Francisco".to_string()),
            zip: Some("94103".to_string()),
            utcoffset: Some(-8),
            lastfix: Some(1234567890),
        });

        let adcom_device = convert_device(&device, &session, geo.as_ref()).unwrap();

        assert_eq!(
            adcom_device.r#type,
            Some(adcom::enums::DeviceType::Phone as i32)
        );
        assert_eq!(adcom_device.ua, Some("Mozilla/5.0".to_string()));
        assert_eq!(adcom_device.make, Some("Apple".to_string()));
        assert_eq!(adcom_device.model, Some("iPhone".to_string()));
        assert_eq!(
            adcom_device.os,
            Some(adcom::enums::OperatingSystem::Ios as i32)
        );
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
        assert_eq!(adcom_device.geo.as_ref().unwrap().lat, Some(37.7749));
        assert_eq!(adcom_device.geo.as_ref().unwrap().lon, Some(-122.4194));
        assert_eq!(adcom_device.geo.as_ref().unwrap().accur, Some(10));
        assert_eq!(
            adcom_device.geo.as_ref().unwrap().country,
            Some("US".to_string())
        );
        assert_eq!(
            adcom_device.geo.as_ref().unwrap().city,
            Some("San Francisco".to_string())
        );
        assert_eq!(
            adcom_device.geo.as_ref().unwrap().zip,
            Some("94103".to_string())
        );
        assert_eq!(adcom_device.geo.as_ref().unwrap().utcoffset, Some(-8));
        assert_eq!(adcom_device.geo.as_ref().unwrap().lastfix, Some(1234567890));

        let device_ext = adcom_device
            .extension_set
            .extension_data(mediation::DEVICE_EXT)
            .unwrap();

        assert_eq!(device_ext.id, Some(id.to_string()));
        assert_eq!(device_ext.launch_ts, Some(1234567890));
        assert_eq!(device_ext.ram_used, Some(1024));
        assert_eq!(device_ext.launch_monotonic_ts, Some(1234567890));
        assert_eq!(device_ext.start_ts, Some(1234567890));
        assert_eq!(device_ext.start_monotonic_ts, Some(1234567890));
        assert_eq!(device_ext.ts, Some(1234567890));
        assert_eq!(device_ext.monotonic_ts, Some(1234567890));
        assert_eq!(device_ext.memory_warnings_ts, vec![1234567890]);
        assert_eq!(device_ext.memory_warnings_monotonic_ts, vec![1234567890]);
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
                format: sdk::AdFormat::Banner,
            }),
            interstitial: Some(json!({"interstitial": "value"})),
            rewarded: Some(json!("rewarded".to_string())),
            auction_pricefloor: 1.0,
        };

        let item = convert_ad_object_to_item(&ad_object).unwrap();

        let mut registry = ExtensionRegistry::new();
        registry.register(mediation::PLACEMENT_EXT);

        let placement = adcom::placement::Placement::decode_with_extensions(
            &mut Cursor::new(item.spec.unwrap()),
            &registry,
        )
        .unwrap();

        let placement_ext = placement
            .extension_set
            .extension_data(mediation::PLACEMENT_EXT)
            .unwrap();

        assert_eq!(item.id, Some("auction_id".to_string()));
        assert_eq!(item.flr, Some(1.0));
        assert_eq!(item.flrcur, Some("USD".to_string()));

        assert_eq!(placement_ext.auction_id, Some("auction_id".to_string()));
        assert_eq!(placement_ext.auction_key, Some("auction_key".to_string()));
        assert_eq!(placement_ext.auction_configuration_id, Some(123));
        assert_eq!(
            placement_ext.auction_configuration_uid,
            Some("auction_configuration_uid".to_string())
        );
        assert_eq!(
            placement_ext.orientation,
            Some(mediation::Orientation::Portrait as i32)
        );
        assert_eq!(
            placement_ext.demands.get("demand_key").unwrap().token,
            Some("token_value".to_string())
        );
        assert_eq!(
            placement_ext.demands.get("demand_key").unwrap().status,
            Some("status_value".to_string())
        );
        assert_eq!(
            placement_ext
                .demands
                .get("demand_key")
                .unwrap()
                .token_finish_ts,
            Some(1234567890)
        );
        assert_eq!(
            placement_ext
                .demands
                .get("demand_key")
                .unwrap()
                .token_start_ts,
            Some(1234567990)
        );
        assert_eq!(
            placement_ext.banner.as_ref().unwrap().format,
            Some(mediation::AdFormat::Banner as i32)
        );
        assert_eq!(
            placement_ext.interstitial,
            Some("{\"interstitial\":\"value\"}".to_string())
        );
        assert_eq!(placement_ext.rewarded, Some("\"rewarded\"".to_string()));
    }
}
