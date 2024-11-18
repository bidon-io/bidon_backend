use crate::com::iabtechlab::adcom::v1 as adcom;
use crate::com::iabtechlab::adcom::v1::context::DistributionChannel;
use crate::com::iabtechlab::openrtb::v3 as openrtb;
use crate::com::iabtechlab::openrtb::v3::AuctionType;
use crate::controllers::adapter::adcom::enums::OperatingSystem;
use crate::galaxy::v1::bidon::{BidonSession, BIDON_AD_OBJECT, BIDON_SESSION};
use crate::galaxy::v1::bidon::{BIDON, BIDON_APP, BIDON_REGS};
use crate::models::{AdFormat, AuctionRequest, DeviceConnectionType, DeviceType, Geo, Segment};
use crate::{galaxy, models};
use galaxy::v1::bidon::{Demand, Orientation};
use models::AdObjectOrientation;
use prost::{Extendable, Message};
use serde_json::Value;
use std::collections::HashMap;
use std::convert::TryFrom;
use std::error::Error;
use std::fmt::format;

//TODO As it takes auction_request's ownership, it should be possible to remove most of the .clone() calls.
pub(crate) fn try_from(
    auction_request: AuctionRequest,
) -> Result<openrtb::Openrtb, Box<dyn Error>> {
    // Convert AuctionRequest to Openrtb::Request
    let mut request = openrtb::Request {
        id: auction_request.ad_object.auction_id.to_owned(),
        test: auction_request.test,
        tmax: auction_request.tmax.map(|t| t as u32),
        at: Some(AuctionType::FirstPrice as i32),
        cur: vec![],  //TODO
        seat: vec![], //TODO
        wseat: None,  //TODO
        context: Some(serialize_context(&auction_request)?),
        source: None, //TODO: Map 'source' if necessary
        item: vec![convert_ad_object_to_item(&auction_request.ad_object)?],
        cdata: None,   //TODO
        package: None, //TODO
        extension_set: Default::default(),
    };

    let bidon_session = convert_session(&auction_request.session);
    request.set_extension_data(BIDON_SESSION, bidon_session)?;

    // Create Openrtb instance with the converted request
    Ok(openrtb::Openrtb {
        ver: Some("3.0".to_string()),                // Set the version as needed
        domainspec: Some("domain_spec".to_string()), // Set the domain spec as needed
        domainver: Some("domain_version".to_string()), // Set the domain version as needed
        payload_oneof: Some(openrtb::openrtb::PayloadOneof::Request(request)),
    })
}

fn serialize_context(auction_request: &AuctionRequest) -> Result<Vec<u8>, Box<dyn Error>> {
    // Create the AdCOM Context message
    let mut context = galaxy::v1::context::Context {
        distribution_channel: DistributionChannel {
            id: None,      // TODO
            name: None,    // TODO
            r#pub: None,   // TODO
            content: None, // TODO
            channel_oneof: Some(adcom::context::distribution_channel::ChannelOneof::App(
                convert_app(&auction_request.app)?,
            )),
        }
            .into(),
        device: convert_device(&auction_request.device, &auction_request.geo).into(),
        user: convert_user(&auction_request.user, &auction_request.segment)?.into(),
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
    api_app: &models::App,
) -> Result<adcom::context::distribution_channel::App, Box<dyn Error>> {
    let mut app = adcom::context::distribution_channel::App {
        ver: api_app.version.clone().into(),
        keywords: None,
        paid: None,
        bundle: api_app.bundle.clone().into(),
        ..Default::default()
    };

    let bidon_app = galaxy::v1::bidon::BidonApp {
        key: api_app.key.clone().into(),
        framework: api_app.framework_version.clone(),
        framework_version: api_app.framework_version.clone(),
        plugin_version: api_app.plugin_version.clone(),
        sdk_version: api_app.sdk_version.clone(),
        skadn: api_app.skadn.clone().unwrap_or(vec![]),
    };

    app.set_extension_data(BIDON_APP, bidon_app)?;
    Ok(app)
}

fn convert_device(api_device: &models::Device, geo: &Option<Geo>) -> adcom::context::Device {
    let mut device = adcom::context::Device {
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
        _ => adcom::enums::ConnectionType::Unknown,
    }
}

fn convert_geo(geo: &models::Geo) -> adcom::context::Geo {
    let geo = adcom::context::Geo {
        r#type: Some(adcom::enums::LocationType::Unknown as i32), // TODO
        lat: geo.lat.map(|t| t as f32),
        lon: geo.lon.map(|t| t as f32),
        accur: geo.accuracy.map(|t| (t as i32)), // TODO check accuracy conversion.
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
    api_user: &models::User,
    segment: &Option<Segment>,
) -> Result<adcom::context::User, Box<dyn Error>> {
    let mut user = adcom::context::User {
        id: api_user.idg.map(|uuid| uuid.to_string()),
        ..Default::default()
    };

    let bidon_user = galaxy::v1::bidon::BidonUser {
        idfa: api_user.idfa.map(|uuid| uuid.to_string()),
        tracking_authorization_status: Some(api_user.tracking_authorization_status.clone()),
        idfv: api_user.idfv.map(|uuid| uuid.to_string()),
        idg: api_user.idg.map(|uuid| uuid.to_string()),
        consent: Some(serde_json::to_string(&api_user.consent)?), // TODO there is a consent field in adcom::User. Should we have it here?
        segments: segment.into_iter().map(convert_segment).collect(),
    };

    user.set_extension_data(BIDON, bidon_user)?;

    Ok(user)
}

fn convert_segment(api_segment: &Segment) -> galaxy::v1::bidon::Segment {
    let mut segment = galaxy::v1::bidon::Segment {
        id: api_segment.id.clone(),
        uid: api_segment.uid.clone(),
        ext: api_segment.ext.clone(),
    };

    segment
}

fn convert_session(api_session: &models::Session) -> BidonSession {
    BidonSession {
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

fn convert_regs(api_regs: &models::Regulations) -> Result<adcom::context::Regs, Box<dyn Error>> {
    let mut regs = adcom::context::Regs {
        coppa: api_regs.coppa,
        gdpr: api_regs.gdpr.clone(),
        ..Default::default()
    };

    let bidon_regs = galaxy::v1::bidon::BidonRegs {
        us_privacy: api_regs.us_privacy.clone(),
        eu_privacy: api_regs.eu_privacy.clone(),
        iab: match api_regs.iab.as_ref() {
            Some(t) => Some(convert_iab(t)?),
            None => None,
        },
    };

    regs.set_extension_data(BIDON_REGS, bidon_regs)?;

    Ok(regs)
}

// Placeholder for IAB conversion
fn convert_iab(iab_json: &HashMap<String, Value>) -> Result<String, Box<dyn Error>> {
    serde_json::to_string(&iab_json).map_err(Into::into)
}

fn convert_ad_object_to_item(
    ad_object: &models::AdObject,
) -> Result<openrtb::Item, Box<dyn Error>> {
    let mut item = openrtb::Item {
        id: ad_object.auction_id.clone(),
        flr: Some(ad_object.auction_pricefloor as f32),
        flrcur: Some("USD".to_string()),
        ..Default::default()
    };

    // Create the AdObjectExtension
    let bidon_ad_object = galaxy::v1::bidon::BidonAdObject {
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

    item.set_extension_data(BIDON_AD_OBJECT, bidon_ad_object)?;

    Ok(item)
}

fn convert_ad_orientation(orientation: &AdObjectOrientation) -> Orientation {
    match orientation {
        AdObjectOrientation::Portrait => Orientation::Portrait,
        AdObjectOrientation::Landscape => Orientation::Landscape,
    }
}

fn convert_banner_ad(api_banner: &models::BannerAdObject) -> galaxy::v1::bidon::BannerAd {
    let mut banner = galaxy::v1::bidon::BannerAd {
        format: Some(convert_ad_format(api_banner.format) as i32),
    };

    banner
}

fn convert_ad_format(format: AdFormat) -> galaxy::v1::bidon::AdFormat {
    match format {
        AdFormat::Banner => galaxy::v1::bidon::AdFormat::Banner,
        AdFormat::Leaderboard => galaxy::v1::bidon::AdFormat::Leaderboard,
        AdFormat::Mrec => galaxy::v1::bidon::AdFormat::Mrec,
        AdFormat::Adaptive => galaxy::v1::bidon::AdFormat::Adaptive,
    }
}

fn convert_demand(
    api_demand: &HashMap<String, Value>,
) -> Result<HashMap<String, Demand>, Box<dyn Error>> {
    let mut demands = HashMap::new();

    for (key, value) in api_demand {
        let map = value
            .as_object()
            .ok_or(format!("Demand value is not an object: {}", value))?;
        let demand = Demand {
            // Assuming Demand has fields that need to be populated from the value
            // Add the necessary field mappings here
            token: match map.get("token") {
                Some(v) => Some(
                    v.as_str()
                        .ok_or(format!(
                            "Token is not a string. Key: {}, value: {}",
                            key, value
                        ))?
                        .to_string(),
                ),
                None => None,
            },
            status: match map.get("status") {
                Some(v) => Some(
                    v.as_str()
                        .ok_or(format!(
                            "Status is not a string. Key: {}, value: {}",
                            key, value
                        ))?
                        .to_string(),
                ),
                None => None,
            },
            token_finish_ts: match map.get("token_finish_ts") {
                Some(v) => Some(
                    v.as_i64()
                        .ok_or(format!(
                            "token_finish_ts is not a number. Key: {}, value: {}",
                            key, value
                        ))?
                ),
                None => None,
            },
            token_start_ts: match map.get("token_start_ts") {
                Some(v) => Some(
                    v.as_i64()
                        .ok_or(format!(
                            "token_start_ts is not a number. Key: {}, value: {}",
                            key, value
                        ))?
                ),
                None => None,
            },
        };
        demands.insert(key.clone(), demand);
    }
    Ok(demands)
}
