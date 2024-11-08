use crate::com::iabtechlab::adcom::v1 as adcom;
use crate::com::iabtechlab::openrtb::v3::Openrtb;
use crate::com::iabtechlab::openrtb::v3 as openrtb;
use crate::models::{AuctionRequest, DeviceConnectionType, DeviceType, Regulations};
use std::convert::TryFrom;
use std::error::Error;
use prost::{EncodeError, Extendable, Message};
use crate::{galaxy, models};
use crate::com::iabtechlab::adcom::v1::DistributionChannel;

pub(crate) fn try_from(auction_request: AuctionRequest) -> Result<Openrtb, dyn Error> {

    // Convert AuctionRequest to Openrtb::Request
    let mut request = openrtb::Request {
        id: auction_request.ad_object.auction_id.to_owned(),
        test: auction_request.test,
        tmax: auction_request.tmax.into(),
        at: None, //TODO
        cur: vec![], //TODO
        seat: vec![], //TODO
        wseat: None, //TODO
        context: Option::from(serialize_context(&auction_request)?),
        source: None, //TODO: Map 'source' if necessary
        item: vec![], //TODO
        cdata: None, //TODO
        package: None, //TODO
    };

    // Create Openrtb instance with the converted request
    Ok(Openrtb {
        ver: Some("3.0".to_string()), // Set the version as needed
        domainspec: Some("domain_spec".to_string()), // Set the domain spec as needed
        domainver: Some("domain_version".to_string()), // Set the domain version as needed
        payload_oneof: Some(openrtb::openrtb::PayloadOneof::Request(request)),
    })
}

fn serialize_context(auction_request: &AuctionRequest) -> Result<Vec<u8>, EncodeError> {
    // Create the AdCOM Context message
    let mut context = galaxy::v1::Context {
        distribution_channel: DistributionChannel {
            id: None, // TODO
            name: None, // TODO
            r#pub: None, // TODO
            content: None, // TODO
            channel_oneof:
            Some(adcom::distribution_channel::ChannelOneof::App(convert_app(&auction_request.app))),
        }.into(),
        device: convert_device(&auction_request.device).into(),
        user: convert_user(&auction_request.user).into(),
        regs: convert_regs(&auction_request.regs).into(),
        restrictions: None, // TODO
    };

    // Serialize the Context message into bytes
    let mut context_bytes = Vec::new();
    context.encode(&mut context_bytes)?;
    Ok(context_bytes)
}

fn convert_app(api_app: &models::App) -> adcom::distribution_channel::App {
    let mut app = adcom::distribution_channel::App {

        // Map standard fields
        domain: None,
        cat: vec![],
        sectcat: vec![],
        pagecat: vec![],
        cattax: None,
        privpolicy: None,
        // app.set_bundle(api_app.bundle);
        storeid: None,
        storeurl: None,
        ver: None,
        // app.set_ver(api_app.version);
        //
        // // Map additional fields
        // app.set_key(api_app.key);
        // app.set_framework(api_app.framework);
        // app.set_framework_version(api_app.framework_version);
        // app.set_plugin_version(api_app.plugin_version);
        // app.set_sdk_version(api_app.sdk_version);
        // app.set_skadn(RepeatedField::from_vec(api_app.skadn,
        // ));
        keywords: None,
        paid: None,
        bundle: None,
        extension_set: Default::default(),
    };
    app.set_extension_data(api_app.extension_data.clone());
    app
}

fn convert_device(api_device: &models::Device) -> adcom::Device {
    let mut device = adcom::Device {
        // Map standard fields
        r#type: convert_device_type(api_device.r#type.clone()).into(),
        ua: api_device.ua.clone().into(),
        ifa: None, // TODO
        dnt: None, // TODO
        lmt: None, // TODO
        make: api_device.make.clone().into(),
        model: api_device.model.clone().into(),
        os: api_device.os.clone().into(),
        osv: api_device.osv.clone().into(),
        hwv: api_device.hwv.clone().into(),
        h: api_device.h.into(),
        w: api_device.w.into(),
        ppi: api_device.ppi.into(),
        pxratio: api_device.pxratio.into(),
        js: Some(api_device.js != 0),
        ip: None, // TODO
        ipv6: None, // TODO
        xff: None, // TODO
        lang: api_device.language.clone().into(),
        carrier: api_device.clone().carrier,
        mccmnc: api_device.clone().mccmnc,
        mccmncsim: None, // TODO
        geofetch: None, // TODO
        contype: Option::from(convert_connection_type(api_device.connection_type).into()),

        // Map additional fields
        // `type`:
        iptr: None,
        geo: None,
        extension_set: Default::default(),
    };
    device
}

fn convert_device_type(device_type: Option<DeviceType>) -> Option<adcom::DeviceType> {
    device_type.map(|dt| match dt {
        DeviceType::Phone => adcom::DeviceType::Phone,
        DeviceType::Tablet => adcom::DeviceType::Tablet,
        _ => adcom::DeviceType::Unknown,
    }).map(Into::into)
}

fn convert_connection_type(connection_type: DeviceConnectionType) -> adcom::ConnectionType {
    match connection_type {
        DeviceConnectionType::Ethernet => adcom::ConnectionType::Wired,
        DeviceConnectionType::Wifi => adcom::ConnectionType::Wifi,
        DeviceConnectionType::CellularUnknown=> adcom::ConnectionType::CellUnknown,
        DeviceConnectionType::Cellular=> adcom::ConnectionType::CellUnknown,
        DeviceConnectionType::Cellular2G => adcom::ConnectionType::Cell2g,
        DeviceConnectionType::Cellular3G => adcom::ConnectionType::Cell3g,
        DeviceConnectionType::Cellular4G => adcom::ConnectionType::Cell4g,
        DeviceConnectionType::Cellular5G => adcom::ConnectionType::Cell5g,
        _ => adcom::ConnectionType::Unknown,
    }
}

// fn convert_user(api_user: &models::User) -> adcom::User {
    // let mut user = adcom::User::new();
    //
    // // Map standard fields
    // user.set_id(api_user.idg);
    //
    // // Map additional fields
    // user.set_idfa(api_user.idfa);
    // user.set_tracking_authorization_status(api_user.tracking_authorization_status);
    // user.set_idfv(api_user.idfv);
    //
    // // Map consent (if structure is known)
    // // If 'consent' is a JSON object, you may need to convert it accordingly
    // if let Some(consent) = api_user.consent {
    //     let consent_struct = convert_consent(consent);
    //     user.set_consent(consent_struct);
    // }

    // // Map 'segments'
    // if let Some(segments_api) = api_user.segments {
    //     let segments = segments_api
    //         .into_iter()
    //         .map(convert_segment)
    //         .collect::<Vec<adcom::Segment>>();
    //     user.set_segments(RepeatedField::from_vec(segments));
    // }
    //
//     user
// }

// fn convert_segment(api_segment: galaxy::bidon::Segment) -> adcom::Segment {
//     let mut segment = galaxy::v1::bidon::Segment{
//         id: None,
//         uid: None,
//         name: None,
//         value: None,
//         ext: None,
//     }
//
//     segment.set_id(api_segment.id);
//     segment.set_uid(api_segment.uid);
//     segment.set_ext(api_segment.ext);
//
//     segment
// }
//
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
