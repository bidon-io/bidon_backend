use crate::com::iabtechlab::adcom::v1 as adcom;
use crate::com::iabtechlab::adcom::v1::context::DistributionChannel;
use crate::com::iabtechlab::openrtb::v3 as openrtb;
use crate::controllers::adapter::adcom::enums::OperatingSystem;
use crate::galaxy::v1::bidon::{BIDON, BIDON_APP};
use crate::models::{AuctionRequest, DeviceConnectionType, DeviceType, Geo, Segment};
use crate::{galaxy, models};
use prost::{EncodeError, Extendable, Message};
use std::convert::TryFrom;
use std::error::Error;

//TODO As it takes auction_request's ownership, it should be possible to remove most of the .clone() calls.
pub(crate) fn try_from(auction_request: AuctionRequest) -> Result<openrtb::Openrtb, Box<dyn Error>> {

    // Convert AuctionRequest to Openrtb::Request
    let mut request = openrtb::Request {
        id: auction_request.ad_object.auction_id.to_owned(),
        test: auction_request.test,
        tmax: auction_request.tmax.map(|t| t as u32),
        at: None, //TODO
        cur: vec![], //TODO
        seat: vec![], //TODO
        wseat: None, //TODO
        context: Option::from(serialize_context(&auction_request)?),
        source: None, //TODO: Map 'source' if necessary
        item: vec![], //TODO
        cdata: None, //TODO
        package: None, //TODO
        extension_set: Default::default(),
    };

    // Create Openrtb instance with the converted request
    Ok(openrtb::Openrtb {
        ver: Some("3.0".to_string()), // Set the version as needed
        domainspec: Some("domain_spec".to_string()), // Set the domain spec as needed
        domainver: Some("domain_version".to_string()), // Set the domain version as needed
        payload_oneof: Some(openrtb::openrtb::PayloadOneof::Request(request)),
    })
}

fn serialize_context(auction_request: &AuctionRequest) -> Result<Vec<u8>, Box<dyn Error>> {
    // Create the AdCOM Context message
    let mut context = galaxy::v1::context::Context {
        distribution_channel: DistributionChannel {
            id: None, // TODO
            name: None, // TODO
            r#pub: None, // TODO
            content: None, // TODO
            channel_oneof:
            Some(adcom::context::distribution_channel::ChannelOneof::App(convert_app(&auction_request.app)?)),
        }.into(),
        device: convert_device(&auction_request.device, &auction_request.geo).into(),
        user: convert_user(&auction_request.user, &auction_request.segment)?.into(),
        regs: None, // convert_regs(&auction_request.regs).into(),
        restrictions: None, // TODO
    };

    // Serialize the Context message into bytes
    let mut context_bytes = Vec::new();
    context.encode(&mut context_bytes)?;
    Ok(context_bytes)
}

fn convert_app(api_app: &models::App) -> Result<adcom::context::distribution_channel::App, Box<dyn Error>> {
    let mut app = adcom::context::distribution_channel::App {

        // Map standard fields
        domain: None,
        cat: vec![],
        sectcat: vec![],
        pagecat: vec![],
        cattax: None,
        privpolicy: None,
        storeid: None,
        storeurl: None,
        ver: api_app.version.clone().into(),

        keywords: None,
        paid: None,
        bundle: api_app.bundle.clone().into(),
        extension_set: Default::default(),
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
        ifa: None, // TODO
        dnt: None, // TODO
        lmt: None, // TODO
        make: api_device.make.clone().into(),
        model: api_device.model.clone().into(),
        os: Option::from(<OperatingSystem as Into<i32>>::into(convert_os(api_device.os.clone()))),
        osv: api_device.osv.clone().into(),
        hwv: api_device.hwv.clone().into(),
        h: api_device.h.into(),
        w: api_device.w.into(),
        ppi: api_device.ppi.into(),
        pxratio: (api_device.pxratio as f32).into(), // TODO validate conversion
        js: Some(api_device.js != 0),
        ip: None, // TODO
        ipv6: None, // TODO
        xff: None, // TODO
        lang: api_device.language.clone().into(),
        carrier: api_device.clone().carrier,
        mccmnc: api_device.clone().mccmnc,
        mccmncsim: None, // TODO
        geofetch: None, // TODO
        contype: Option::from(<adcom::enums::ConnectionType as Into<i32>>::into(convert_connection_type(api_device.connection_type))),

        // Map additional fields
        // `type`:
        iptr: None,
        geo: geo.clone().map(|g| convert_geo(&g)),
        extension_set: Default::default(),
    };
    device
}

fn convert_device_type(device_type: Option<DeviceType>) -> Option<adcom::enums::DeviceType> {
    device_type.map(|dt| match dt {
        DeviceType::Phone => adcom::enums::DeviceType::Phone,
        DeviceType::Tablet => adcom::enums::DeviceType::Tablet,
    }).map(Into::into)
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
        r#type: Option::from(adcom::enums::LocationType::Unknown as i32), // TODO
        lat: geo.lat.map(|t|t as f32),
        lon: geo.lon.map(|t|t as f32),
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

fn convert_user(api_user: &models::User, segment: &Option<Segment>) -> Result<adcom::context::User, Box<dyn Error>>  {
    let mut user = adcom::context::User {
        id: api_user.idg.map(|uuid| uuid.to_string()),
        ..Default::default()
    };

    let bidon_user = galaxy::v1::bidon::BidonUser {
        idfa: api_user.idfa.map(|uuid| uuid.to_string()),
        tracking_authorization_status: Option::from(api_user.tracking_authorization_status.clone()),
        idfv: api_user.idfv.map(|uuid| uuid.to_string()),
        idg: api_user.idg.map(|uuid| uuid.to_string()),
        consent: Option::from(serde_json::to_string(&api_user.consent)?), // TODO there is a consent field in adcom::User. Should we have it here?
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

// // Placeholder for consent conversion
// fn convert_consent(consent_json: serde_json::Value) -> adcom::Consent {
//     let mut consent = adcom::Consent::new();
//     // Map fields from consent_json to consent message
//     // This requires defining the Consent message structure
//     consent
// }
// fn convert_regs(api_regs: &Option<Regulations>) -> adcom::Regs {
//     let mut regs = adcom::Regs::new();
//
//     regs.set_coppa(api_regs.coppa);
//
//     // Map additional fields
//     regs.set_gdpr(api_regs.gdpr);
//     regs.set_us_privacy(api_regs.us_privacy);
//     regs.set_eu_privacy(api_regs.eu_privacy);
//
//     // Map 'iab' if structure is known
//     if let Some(iab) = api_regs.iab {
//         let iab_struct = convert_iab(iab);
//         regs.set_iab(iab_struct);
//     }
//
//     regs
// }
//
// // Placeholder for IAB conversion
// fn convert_iab(iab_json: serde_json::Value) -> adcom::Iab {
//     let mut iab = adcom::Iab::new();
//     // Map fields from iab_json to IAB message
//     // This requires defining the Iab message structure
//     iab
// }
//
// fn convert_ad_object_to_impressions(api_ad_object: OpenApiAdObject) -> RepeatedField<openrtb::Impression> {
//     let mut impressions = Vec::new();
//
//     let mut imp = openrtb::Impression::new();
//
//     imp.set_id("1".to_string()); // Assign a unique ID
//
//     // Map 'banner' ad
//     if let Some(banner_api) = api_ad_object.banner {
//         let banner = convert_banner_ad(banner_api);
//         imp.set_banner(banner);
//     }
//
//     // Map 'interstitial' ad
//     if let Some(interstitial_api) = api_ad_object.interstitial {
//         // Map as a banner or video depending on your implementation
//         // For example, if it's a banner:
//         let banner = convert_interstitial_ad(interstitial_api);
//         imp.set_banner(banner);
//     }
//
//     // Map 'rewarded' ad
//     if let Some(rewarded_api) = api_ad_object.rewarded {
//         // Map as a video ad
//         let video = convert_rewarded_ad(rewarded_api);
//         imp.set_video(video);
//     }
//
//     // Map other fields like 'auction_pricefloor'
//     if let Some(pricefloor) = api_ad_object.auction_pricefloor {
//         imp.set_bidfloor(pricefloor);
//     }
//
//     // Add imp to the list
//     impressions.push(imp);
//
//     RepeatedField::from_vec(impressions)
// }
//
// fn convert_banner_ad(api_banner: OpenApiBannerAdObject) -> adcom::Banner {
//     let mut banner = adcom::Banner::new();
//
//     // Map 'format'
//     banner.set_format(convert_ad_format(api_banner.format));
//
//     // Map other fields as needed
//
//     banner
// }
//
// fn convert_ad_format(format_str: String) -> adcom::AdFormat {
//     match format_str.as_str() {
//         "BANNER" => adcom::Ad::BANNER,
//         "LEADERBOARD" => adcom::AdFormat::LEADERBOARD,
//         "MREC" => adcom::AdFormat::MREC,
//         "ADAPTIVE" => adcom::AdFormat::ADAPTIVE,
//         _ => adcom::AdFormat::UNKNOWN_FORMAT,
//     }
// }
//
// fn convert_interstitial_ad(api_interstitial: OpenApiInterstitialAdObject) -> adcom::Banner {
//     let mut banner = adcom::Banner::new();
//
//     // Map fields specific to interstitial ads
//     // ...
//
//     banner
// }
//
// fn convert_rewarded_ad(api_rewarded: OpenApiRewardedAdObject) -> adcom::Video {
//     let mut video = adcom::Video::new();
//
//     // Map fields specific to rewarded ads
//     // ...
//
//     video
// }
