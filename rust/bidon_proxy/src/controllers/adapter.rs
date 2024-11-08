use crate::com::iabtechlab::adcom::v1 as adcom;
use crate::com::iabtechlab::openrtb::v3::Openrtb;
use crate::com::iabtechlab::openrtb::v3 as openrtb;
use crate::models::AuctionRequest;
use std::convert::TryFrom;
use std::error::Error;
use crate::models;

pub(crate) fn try_from(auction_request: AuctionRequest) -> Result<Openrtb, dyn Error> {

    // Convert AuctionRequest to Openrtb::Request
    let mut request = openrtb::Request {
        id: auction_request.ad_object.auction_id,
        test: auction_request.test,
        tmax: auction_request.tmax,
        context: auction_request.context,
        app: convert_app(auction_request.app),
        device: convert_device(auction_request.device),
        user: convert_user(auction_request.user),
        regs: convert_regs(auction_request.regs),
        source: None, //TODO: Map 'source' if necessary
        imp: convert_ad_object_to_impressions(auction_request.ad_object),
    };

    // Create Openrtb instance with the converted request
    Ok(Openrtb {
        ver: Some("3.0".to_string()), // Set the version as needed
        domainspec: Some("domain_spec".to_string()), // Set the domain spec as needed
        domainver: Some("domain_version".to_string()), // Set the domain version as needed
        request: Some(request),
        response: None,
        ext: None,
    })
}

fn serialize_context(auction_request: &AuctionRequest) -> Vec<u8> {
    // Create the AdCOM Context message
    let mut context = adcom::Context::new();

    // Map 'App'
    let app = convert_app(auction_request.app);
    context.set_app(app);

    // Map 'Device'
    let device = convert_device(auction_request.device);
    context.set_device(device);

    // Map 'User'
    let user = convert_user(auction_request.user);
    context.set_user(user);

    // Map 'Regs'
    if let Some(regs_api) = auction_request.regs {
        let regs = convert_regs(regs_api);
        context.set_regs(regs);
    }

    // Map 'Restrictions' if necessary
    // context.set_restrictions(...);

    // Serialize the Context message into bytes
    let mut context_bytes = Vec::new();
    context.write_to_vec(&mut context_bytes).unwrap();

    context_bytes
}

fn convert_app(api_app: &models::App) -> adcom::App {
    let mut app = adcom::App {}

    // Map standard fields
    app.set_bundle(api_app.bundle);
    app.set_ver(api_app.version);

    // Map additional fields
    app.set_key(api_app.key);
    app.set_framework(api_app.framework);
    app.set_framework_version(api_app.framework_version);
    app.set_plugin_version(api_app.plugin_version);
    app.set_sdk_version(api_app.sdk_version);
    app.set_skadn(RepeatedField::from_vec(api_app.skadn));

    app
}

fn convert_device(api_device: models::Device) -> adcom::Device {
    let mut device = adcom::Device {
        // Map standard fields
        ua: api_device.ua.into(),
        make: api_device.make.into(),
        model: api_device.model.into(),
        os: api_device.os.into(),
        osv: api_device.osv.into(),
        hwv: api_device.hwv.into(),
        h: api_device.h.into(),
        w: api_device.w.into(),
        ppi: api_device.ppi.into(),
        pxratio: api_device.pxratio.into(),
        js: api_device.js.into(),
        language: api_device.language,
        carrier: api_device.carrier,
        mccmnc: api_device.mccmnc,
        connectiontype: convert_connection_type(api_device.connection_type.into()),

        // Map additional fields
        `type`: convert_device_type(api_device.r#type),
    }
    device
}

fn convert_device_type(device_type: String) -> adcom::DeviceType {
    match device_type.as_str() {
        "PHONE" => adcom::DeviceType::PHONE,
        "TABLET" => adcom::DeviceType::TABLET,
        _ => adcom::DeviceType::UNKNOWN_DEVICE_TYPE,
    }
}

fn convert_connection_type(connection_type: String) -> adcom::ConnectionType {
    match connection_type.as_str() {
        "ETHERNET" => adcom::ConnectionType::ETHERNET,
        "WIFI" => adcom::ConnectionType::WIFI,
        "CELLULAR" => adcom::ConnectionType::CELLULAR,
        "CELLULAR_UNKNOWN" => adcom::ConnectionType::CELLULAR_UNKNOWN,
        "CELLULAR_2_G" => adcom::ConnectionType::CELLULAR_2G,
        "CELLULAR_3_G" => adcom::ConnectionType::CELLULAR_3G,
        "CELLULAR_4_G" => adcom::ConnectionType::CELLULAR_4G,
        "CELLULAR_5_G" => adcom::ConnectionType::CELLULAR_5G,
        _ => adcom::ConnectionType::UNKNOWN_CONNECTION_TYPE,
    }
}

fn convert_user(api_user: OpenApiUser) -> adcom::User {
    let mut user = adcom::User::new();

    // Map standard fields
    user.set_id(api_user.idg);

    // Map additional fields
    user.set_idfa(api_user.idfa);
    user.set_tracking_authorization_status(api_user.tracking_authorization_status);
    user.set_idfv(api_user.idfv);

    // Map consent (if structure is known)
    // If 'consent' is a JSON object, you may need to convert it accordingly
    if let Some(consent) = api_user.consent {
        let consent_struct = convert_consent(consent);
        user.set_consent(consent_struct);
    }

    // Map 'segments'
    if let Some(segments_api) = api_user.segments {
        let segments = segments_api
            .into_iter()
            .map(convert_segment)
            .collect::<Vec<adcom::Segment>>();
        user.set_segments(RepeatedField::from_vec(segments));
    }

    user
}

fn convert_segment(api_segment: OpenApiSegment) -> adcom::Segment {
    let mut segment = adcom::Segment::new();

    segment.set_id(api_segment.id);
    segment.set_uid(api_segment.uid);
    segment.set_ext(api_segment.ext);

    segment
}

// Placeholder for consent conversion
fn convert_consent(consent_json: serde_json::Value) -> adcom::Consent {
    let mut consent = adcom::Consent::new();
    // Map fields from consent_json to consent message
    // This requires defining the Consent message structure
    consent
}
fn convert_regs(api_regs: OpenApiRegs) -> adcom::Regs {
    let mut regs = adcom::Regs::new();

    regs.set_coppa(api_regs.coppa);

    // Map additional fields
    regs.set_gdpr(api_regs.gdpr);
    regs.set_us_privacy(api_regs.us_privacy);
    regs.set_eu_privacy(api_regs.eu_privacy);

    // Map 'iab' if structure is known
    if let Some(iab) = api_regs.iab {
        let iab_struct = convert_iab(iab);
        regs.set_iab(iab_struct);
    }

    regs
}

// Placeholder for IAB conversion
fn convert_iab(iab_json: serde_json::Value) -> adcom::Iab {
    let mut iab = adcom::Iab::new();
    // Map fields from iab_json to IAB message
    // This requires defining the Iab message structure
    iab
}

fn convert_ad_object_to_impressions(api_ad_object: OpenApiAdObject) -> RepeatedField<openrtb::Impression> {
    let mut impressions = Vec::new();

    let mut imp = openrtb::Impression::new();

    imp.set_id("1".to_string()); // Assign a unique ID

    // Map 'banner' ad
    if let Some(banner_api) = api_ad_object.banner {
        let banner = convert_banner_ad(banner_api);
        imp.set_banner(banner);
    }

    // Map 'interstitial' ad
    if let Some(interstitial_api) = api_ad_object.interstitial {
        // Map as a banner or video depending on your implementation
        // For example, if it's a banner:
        let banner = convert_interstitial_ad(interstitial_api);
        imp.set_banner(banner);
    }

    // Map 'rewarded' ad
    if let Some(rewarded_api) = api_ad_object.rewarded {
        // Map as a video ad
        let video = convert_rewarded_ad(rewarded_api);
        imp.set_video(video);
    }

    // Map other fields like 'auction_pricefloor'
    if let Some(pricefloor) = api_ad_object.auction_pricefloor {
        imp.set_bidfloor(pricefloor);
    }

    // Add imp to the list
    impressions.push(imp);

    RepeatedField::from_vec(impressions)
}

fn convert_banner_ad(api_banner: OpenApiBannerAdObject) -> adcom::Banner {
    let mut banner = adcom::Banner::new();

    // Map 'format'
    banner.set_format(convert_ad_format(api_banner.format));

    // Map other fields as needed

    banner
}

fn convert_ad_format(format_str: String) -> adcom::AdFormat {
    match format_str.as_str() {
        "BANNER" => adcom::AdFormat::BANNER,
        "LEADERBOARD" => adcom::AdFormat::LEADERBOARD,
        "MREC" => adcom::AdFormat::MREC,
        "ADAPTIVE" => adcom::AdFormat::ADAPTIVE,
        _ => adcom::AdFormat::UNKNOWN_FORMAT,
    }
}

fn convert_interstitial_ad(api_interstitial: OpenApiInterstitialAdObject) -> adcom::Banner {
    let mut banner = adcom::Banner::new();

    // Map fields specific to interstitial ads
    // ...

    banner
}

fn convert_rewarded_ad(api_rewarded: OpenApiRewardedAdObject) -> adcom::Video {
    let mut video = adcom::Video::new();

    // Map fields specific to rewarded ads
    // ...

    video
}
