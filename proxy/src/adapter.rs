use crate::adapter::adcom::enums::OperatingSystem;
use crate::bidon::v1::mediation::{
    MediationSession, MEDIATION_AD_OBJECT, MEDIATION_SESSION, MEDIATION_USER,
};
use crate::bidon::v1::mediation::{MEDIATION_APP, MEDIATION_REGS};
use crate::com::iabtechlab::adcom::v1 as adcom;
use crate::com::iabtechlab::adcom::v1::context::DistributionChannel;
use crate::com::iabtechlab::openrtb::v3 as openrtb;
use crate::com::iabtechlab::openrtb::v3::AuctionType;
use crate::protocol::{AdFormat, AuctionRequest, DeviceConnectionType, DeviceType, Geo, Segment};
use crate::{bidon, protocol};
use anyhow::{anyhow, Result};
use bidon::v1::mediation::{Demand, Orientation};
use prost::{Extendable, Message};
use protocol::AdObjectOrientation;
use serde_json::Value;
use std::collections::HashMap;

pub(crate) fn try_from(
    auction_request: AuctionRequest,
    bidon_version: &String,
) -> Result<openrtb::Openrtb> {
    // Convert AuctionRequest to Openrtb::Request
    let mut request = openrtb::Request {
        id: auction_request.ad_object.auction_id.to_owned(),
        test: auction_request.test,
        tmax: auction_request.tmax.map(|t| t as u32),
        at: Some(AuctionType::FirstPrice as i32),
        context: Some(serialize_context(&auction_request, bidon_version)?),
        item: vec![convert_ad_object_to_item(&auction_request.ad_object)?],
        ..Default::default()
    };

    let bidon_session = convert_session(&auction_request.session);
    request.set_extension_data(MEDIATION_SESSION, bidon_session)?;

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
    let context = bidon::v1::context::Context {
        distribution_channel: DistributionChannel {
            channel_oneof: Some(adcom::context::distribution_channel::ChannelOneof::App(
                convert_app(&auction_request.app, bidon_version)?,
            )),
            ..Default::default()
        }
        .into(),
        device: convert_device(&auction_request.device, auction_request.geo.as_ref()).into(),
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
    api_app: &protocol::App,
    bidon_version: &String,
) -> Result<adcom::context::distribution_channel::App> {
    let mut app = adcom::context::distribution_channel::App {
        ver: api_app.version.clone().into(),
        keywords: None,
        paid: None,
        bundle: api_app.bundle.clone().into(),
        ..Default::default()
    };

    let bidon_app = bidon::v1::mediation::MediationApp {
        key: api_app.key.clone().into(),
        framework: api_app.framework_version.clone(),
        framework_version: api_app.framework_version.clone(),
        plugin_version: api_app.plugin_version.clone(),
        sdk_version: api_app.sdk_version.clone(),
        skadn: api_app.skadn.clone().unwrap_or(vec![]),
        bidon_version: Some(bidon_version.clone()),
    };

    app.set_extension_data(MEDIATION_APP, bidon_app)?;
    Ok(app)
}

fn convert_device(api_device: &protocol::Device, geo: Option<&Geo>) -> adcom::context::Device {
    let device = adcom::context::Device {
        // Map standard fields
        r#type: convert_device_type(api_device.r#type.clone()).map(|dt| dt as i32),
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
    device
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

fn convert_geo(geo: &protocol::Geo) -> adcom::context::Geo {
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
    api_user: &protocol::User,
    segment: Option<&Segment>,
) -> Result<adcom::context::User> {
    let mut user = adcom::context::User {
        id: api_user.idg.map(|uuid| uuid.to_string()),
        ..Default::default()
    };

    let mediation_user = bidon::v1::mediation::MediationUser {
        idfa: api_user.idfa.map(|uuid| uuid.to_string()),
        tracking_authorization_status: Some(api_user.tracking_authorization_status.clone()),
        idfv: api_user.idfv.map(|uuid| uuid.to_string()),
        idg: api_user.idg.map(|uuid| uuid.to_string()),
        consent: Some(serde_json::to_string(&api_user.consent)?), // TODO there is a consent field in adcom::User. Should we have it here?
        segments: segment.into_iter().map(convert_segment).collect(),
    };

    user.set_extension_data(MEDIATION_USER, mediation_user)?;

    Ok(user)
}

fn convert_segment(api_segment: &Segment) -> bidon::v1::mediation::Segment {
    let segment = bidon::v1::mediation::Segment {
        id: api_segment.id.clone(),
        uid: api_segment.uid.clone(),
        ext: api_segment.ext.clone(),
    };

    segment
}

fn convert_session(api_session: &protocol::Session) -> MediationSession {
    MediationSession {
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

fn convert_regs(api_regs: &protocol::Regulations) -> Result<adcom::context::Regs> {
    let mut regs = adcom::context::Regs {
        coppa: api_regs.coppa,
        gdpr: api_regs.gdpr.clone(),
        ..Default::default()
    };

    let mediation_regs = bidon::v1::mediation::MediationRegs {
        us_privacy: api_regs.us_privacy.clone(),
        eu_privacy: api_regs.eu_privacy.clone(),
        iab: match api_regs.iab.as_ref() {
            Some(t) => Some(convert_iab(t)?),
            None => None,
        },
    };

    regs.set_extension_data(MEDIATION_REGS, mediation_regs)?;

    Ok(regs)
}

fn convert_iab(iab_json: &HashMap<String, Value>) -> Result<String> {
    serde_json::to_string(&iab_json).map_err(Into::into)
}

fn convert_ad_object_to_item(ad_object: &protocol::AdObject) -> Result<openrtb::Item> {
    let mut item = openrtb::Item {
        id: ad_object.auction_id.clone(),
        flr: Some(ad_object.auction_pricefloor as f32),
        flrcur: Some("USD".to_string()),
        ..Default::default()
    };

    // Create the AdObjectExtension
    let mediation_ad_object = bidon::v1::mediation::MediationAdObject {
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

    item.set_extension_data(MEDIATION_AD_OBJECT, mediation_ad_object)?;

    Ok(item)
}

fn convert_ad_orientation(orientation: &AdObjectOrientation) -> Orientation {
    match orientation {
        AdObjectOrientation::Portrait => Orientation::Portrait,
        AdObjectOrientation::Landscape => Orientation::Landscape,
    }
}

fn convert_banner_ad(api_banner: &protocol::BannerAdObject) -> bidon::v1::mediation::BannerAd {
    let banner = bidon::v1::mediation::BannerAd {
        format: Some(convert_ad_format(api_banner.format) as i32),
    };

    banner
}

fn convert_ad_format(format: AdFormat) -> bidon::v1::mediation::AdFormat {
    match format {
        AdFormat::Banner => bidon::v1::mediation::AdFormat::Banner,
        AdFormat::Leaderboard => bidon::v1::mediation::AdFormat::Leaderboard,
        AdFormat::Mrec => bidon::v1::mediation::AdFormat::Mrec,
        AdFormat::Adaptive => bidon::v1::mediation::AdFormat::Adaptive,
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

#[cfg(test)]
mod tests {
    use super::*;
    use crate::bidon::v1::mediation::MEDIATION_USER;
    use crate::protocol::{App, Device, Session, User};
    use serde_json::json;
    use std::collections::HashMap;
    use uuid::Uuid;

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
    fn test_convert_device() {
        let api_device = Device {
            r#type: Some(DeviceType::Phone),
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
        };

        let geo = Some(Geo {
            lat: Some(37.7749),
            lon: Some(-122.4194),
            accuracy: Some(10.6),
            country: Some("US".to_string()),
            city: Some("San Francisco".to_string()),
            zip: Some("94103".to_string()),
            utcoffset: Some(-8),
            lastfix: Some(1234567890),
        });

        let device = convert_device(&api_device, geo.as_ref());

        assert_eq!(device.r#type, Some(adcom::enums::DeviceType::Phone as i32));
        assert_eq!(device.ua, Some("Mozilla/5.0".to_string()));
        assert_eq!(device.make, Some("Apple".to_string()));
        assert_eq!(device.model, Some("iPhone".to_string()));
        assert_eq!(device.os, Some(adcom::enums::OperatingSystem::Ios as i32));
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
        assert_eq!(device.geo.as_ref().unwrap().lat, Some(37.7749));
        assert_eq!(device.geo.as_ref().unwrap().lon, Some(-122.4194));
        assert_eq!(device.geo.as_ref().unwrap().accur, Some(10));
        assert_eq!(device.geo.as_ref().unwrap().country, Some("US".to_string()));
        assert_eq!(
            device.geo.as_ref().unwrap().city,
            Some("San Francisco".to_string())
        );
        assert_eq!(device.geo.as_ref().unwrap().zip, Some("94103".to_string()));
        assert_eq!(device.geo.as_ref().unwrap().utcoffset, Some(-8));
        assert_eq!(device.geo.as_ref().unwrap().lastfix, Some(1234567890));
    }

    #[test]
    fn test_convert_user() {
        let idfa = Uuid::new_v4();
        let idfv = Uuid::new_v4();
        let idg = Uuid::new_v4();

        let api_user = User {
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

        let segment = Some(Segment {
            id: Some("segment_id".to_string()),
            uid: Some("segment_uid".to_string()),
            ext: None,
        });

        let user = convert_user(&api_user, segment.as_ref()).unwrap();

        assert_eq!(user.id, Some(idg.to_string()));
        let bidon_user = user.extension_set.extension_data(MEDIATION_USER).unwrap();
        assert_eq!(bidon_user.idfa, Some(idfa.to_string()));
        assert_eq!(bidon_user.idfv, Some(idfv.to_string()));
        assert_eq!(
            bidon_user.tracking_authorization_status,
            Some("authorized".to_string())
        );
        assert_eq!(
            bidon_user.consent,
            Some("{\"meta\":{\"consent\":true}}".to_string())
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
        let bidon_app = app.extension_set.extension_data(MEDIATION_APP).unwrap();
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
        let api_regs = protocol::Regulations {
            coppa: Some(true),
            gdpr: Some(true),
            us_privacy: Some("1YNN".to_string()),
            eu_privacy: Some("1".to_string()),
            iab: Some(HashMap::from([("key".to_string(), json!("value"))])),
        };

        let regs = convert_regs(&api_regs).unwrap();

        assert_eq!(regs.coppa, Some(true));
        assert_eq!(regs.gdpr, Some(true));
        let bidon_regs = regs.extension_set.extension_data(MEDIATION_REGS).unwrap();
        assert_eq!(bidon_regs.us_privacy, Some("1YNN".to_string()));
        assert_eq!(bidon_regs.eu_privacy, Some("1".to_string()));
        assert_eq!(bidon_regs.iab, Some("{\"key\":\"value\"}".to_string()));
    }

    #[test]
    fn test_convert_ad_object_to_item() {
        let ad_object = protocol::AdObject {
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
            banner: Some(protocol::BannerAdObject {
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
        let bidon_ad_object = item
            .extension_set
            .extension_data(MEDIATION_AD_OBJECT)
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
            Some(bidon::v1::mediation::AdFormat::Banner as i32)
        );
        assert_eq!(
            bidon_ad_object.interstitial,
            Some("{\"interstitial\":\"value\"}".to_string())
        );
        assert_eq!(bidon_ad_object.rewarded, Some("\"rewarded\"".to_string()));
    }
}
