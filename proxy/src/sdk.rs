#![allow(unused_qualifications)]

use crate::sdk;

/// Format of the banner ad
/// Enumeration of values.
/// Since this enum's variants do not hold data, we can easily define them as `#[repr(C)]`
/// which helps with FFI.
#[allow(non_camel_case_types)]
#[repr(C)]
#[derive(
    Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, serde::Serialize, serde::Deserialize, Hash,
)]
pub enum AdFormat {
    #[serde(rename = "BANNER")]
    Banner,
    #[serde(rename = "LEADERBOARD")]
    Leaderboard,
    #[serde(rename = "MREC")]
    Mrec,
    #[serde(rename = "ADAPTIVE")]
    Adaptive,
}

impl std::str::FromStr for AdFormat {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        match s {
            "BANNER" => std::result::Result::Ok(AdFormat::Banner),
            "LEADERBOARD" => std::result::Result::Ok(AdFormat::Leaderboard),
            "MREC" => std::result::Result::Ok(AdFormat::Mrec),
            "ADAPTIVE" => std::result::Result::Ok(AdFormat::Adaptive),
            _ => std::result::Result::Err(format!("Value not valid: {}", s)),
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct AdObject {
    /// ID of the bidding configuration
    #[serde(rename = "auction_configuration_id")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub auction_configuration_id: Option<i64>,

    /// UID of the bidding configuration
    #[serde(rename = "auction_configuration_uid")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub auction_configuration_uid: Option<String>,

    /// Unique identifier for the bidding
    #[serde(rename = "auction_id")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub auction_id: Option<String>,

    /// Generated key for the bidding request
    #[serde(rename = "auction_key")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub auction_key: Option<String>,

    /// PriceFloor for the bidding
    #[serde(rename = "auction_pricefloor")]
    #[validate(range(min = 0))]
    pub auction_pricefloor: f64,

    #[serde(rename = "banner")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub banner: Option<sdk::BannerAdObject>,

    /// Map of demands
    #[serde(rename = "demands")]
    pub demands: std::collections::HashMap<String, serde_json::Value>,

    /// Empty schema for interstitial ad configuration
    #[serde(rename = "interstitial")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub interstitial: Option<serde_json::Value>,

    #[serde(rename = "orientation")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub orientation: Option<sdk::AdObjectOrientation>,

    /// Empty schema for rewarded ad configuration
    #[serde(rename = "rewarded")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub rewarded: Option<serde_json::Value>,
}

impl AdObject {
    #[allow(clippy::new_without_default)]
    pub fn new(
        auction_pricefloor: f64,
        demands: std::collections::HashMap<String, serde_json::Value>,
    ) -> AdObject {
        AdObject {
            auction_configuration_id: None,
            auction_configuration_uid: None,
            auction_id: None,
            auction_key: None,
            auction_pricefloor,
            banner: None,
            demands,
            interstitial: None,
            orientation: None,
            rewarded: None,
        }
    }
}

/// Orientation of the ad
/// Enumeration of values.
/// Since this enum's variants do not hold data, we can easily define them as `#[repr(C)]`
/// which helps with FFI.
#[allow(non_camel_case_types)]
#[repr(C)]
#[derive(
    Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, serde::Serialize, serde::Deserialize, Hash,
)]
pub enum AdObjectOrientation {
    #[serde(rename = "PORTRAIT")]
    Portrait,
    #[serde(rename = "LANDSCAPE")]
    Landscape,
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct AdUnit {
    /// Type of bid associated with the ad unit
    #[serde(rename = "bid_type")]
    pub bid_type: String,

    /// Identifier for the demand source
    #[serde(rename = "demand_id")]
    pub demand_id: String,

    /// Additional properties for the ad unit
    #[serde(rename = "ext")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ext: Option<std::collections::HashMap<String, serde_json::Value>>,

    /// Label for the ad unit
    #[serde(rename = "label")]
    pub label: String,

    /// Optional price floor for the ad unit
    #[serde(rename = "pricefloor")]
    #[validate(range(min = 0))]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub pricefloor: Option<f64>,

    /// Unique identifier for the ad unit
    #[serde(rename = "uid")]
    pub uid: String,
}

impl AdUnit {
    #[allow(clippy::new_without_default)]
    pub fn new(bid_type: String, demand_id: String, label: String, uid: String) -> AdUnit {
        AdUnit {
            bid_type,
            demand_id,
            ext: None,
            label,
            pricefloor: None,
            uid,
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct Adapter {
    #[serde(rename = "sdk_version")]
    pub sdk_version: String,

    #[serde(rename = "version")]
    pub version: String,
}

impl Adapter {
    #[allow(clippy::new_without_default)]
    pub fn new(sdk_version: String, version: String) -> Adapter {
        Adapter {
            sdk_version,
            version,
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct App {
    /// The bundle identifier of the application, typically in reverse domain name notation (e.g., com.example.myapp).
    #[serde(rename = "bundle")]
    pub bundle: String,

    /// The name of the framework used by the application (e.g., React Native, Flutter, etc.).
    #[serde(rename = "framework")]
    pub framework: String,

    /// The version of the framework used by the application, specifying compatibility.
    #[serde(rename = "framework_version")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub framework_version: Option<String>,

    /// A unique key or identifier for the application.
    #[serde(rename = "key")]
    pub key: String,

    /// The version of the plugin integrated into the application
    #[serde(rename = "plugin_version")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub plugin_version: Option<String>,

    /// The version of the SDK used in the application.
    #[serde(rename = "sdk_version")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub sdk_version: Option<String>,

    /// An array of SKAdNetwork IDs for ad attribution, used primarily for iOS applications.
    #[serde(rename = "skadn")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub skadn: Option<Vec<String>>,

    /// The version of the application, typically following semantic versioning (e.g., 1.0.0).
    #[serde(rename = "version")]
    pub version: String,
}

impl App {
    #[allow(clippy::new_without_default)]
    pub fn new(bundle: String, framework: String, key: String, version: String) -> App {
        App {
            bundle,
            framework,
            framework_version: None,
            key,
            plugin_version: None,
            sdk_version: None,
            skadn: None,
            version,
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct AuctionAdUnitResult {
    /// Label of the ad unit
    #[serde(rename = "ad_unit_label")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ad_unit_label: Option<String>,

    /// UID of the ad unit
    #[serde(rename = "ad_unit_uid")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ad_unit_uid: Option<String>,

    #[serde(rename = "bid_type")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub bid_type: Option<serde_json::Value>,

    /// ID of the demand source for the ad unit
    #[serde(rename = "demand_id")]
    pub demand_id: String,

    /// Error message associated with the ad unit, if applicable
    #[serde(rename = "error_message")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub error_message: Option<String>,

    /// Timestamp when the ad fill finished
    #[serde(rename = "fill_finish_ts")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub fill_finish_ts: Option<i64>,

    /// Timestamp when the ad fill started
    #[serde(rename = "fill_start_ts")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub fill_start_ts: Option<i64>,

    /// Price associated with the ad unit
    #[serde(rename = "price")]
    #[validate(range(min = 0))]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub price: Option<f64>,

    #[serde(rename = "status")]
    pub status: sdk::AuctionAdUnitResultStatus,

    /// Timestamp when the token process finished
    #[serde(rename = "token_finish_ts")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub token_finish_ts: Option<i64>,

    /// Timestamp when the token process started
    #[serde(rename = "token_start_ts")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub token_start_ts: Option<i64>,
}

impl AuctionAdUnitResult {
    #[allow(clippy::new_without_default)]
    pub fn new(demand_id: String, status: sdk::AuctionAdUnitResultStatus) -> AuctionAdUnitResult {
        AuctionAdUnitResult {
            ad_unit_label: None,
            ad_unit_uid: None,
            bid_type: None,
            demand_id,
            error_message: None,
            fill_finish_ts: None,
            fill_start_ts: None,
            price: None,
            status,
            token_finish_ts: None,
            token_start_ts: None,
        }
    }
}

/// Status of the ad unit
/// Enumeration of values.
/// Since this enum's variants do not hold data, we can easily define them as `#[repr(C)]`
/// which helps with FFI.
#[allow(non_camel_case_types)]
#[repr(C)]
#[derive(
    Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, serde::Serialize, serde::Deserialize, Hash,
)]
pub enum AuctionAdUnitResultStatus {
    #[serde(rename = "")]
    Empty,
    #[serde(rename = "SUCCESS")]
    Success,
    #[serde(rename = "FAIL")]
    Fail,
    #[serde(rename = "PENDING")]
    Pending,
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct AuctionRequest {
    #[serde(rename = "app")]
    pub app: sdk::App,

    #[serde(rename = "device")]
    pub device: sdk::Device,

    #[serde(rename = "ext")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ext: Option<String>,

    #[serde(rename = "geo")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub geo: Option<sdk::Geo>,

    #[serde(rename = "regs")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub regs: Option<sdk::Regulations>,

    #[serde(rename = "segment")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub segment: Option<sdk::Segment>,

    #[serde(rename = "session")]
    pub session: sdk::Session,

    #[serde(rename = "token")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub token: Option<String>,

    #[serde(rename = "user")]
    pub user: sdk::User,

    #[serde(rename = "ad_object")]
    pub ad_object: sdk::AdObject,

    #[serde(rename = "adapters")]
    pub adapters: std::collections::HashMap<String, sdk::Adapter>,

    /// Flag indicating that the request is a test
    #[serde(rename = "test")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub test: Option<bool>,

    /// Maximum response time for the server before timeout
    #[serde(rename = "tmax")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub tmax: Option<i64>,
}

impl AuctionRequest {
    #[allow(clippy::new_without_default)]
    pub fn new(
        app: sdk::App,
        device: sdk::Device,
        session: sdk::Session,
        user: sdk::User,
        ad_object: sdk::AdObject,
        adapters: std::collections::HashMap<String, sdk::Adapter>,
    ) -> AuctionRequest {
        AuctionRequest {
            app,
            device,
            ext: None,
            geo: None,
            regs: None,
            segment: None,
            session,
            token: None,
            user,
            ad_object,
            adapters,
            test: None,
            tmax: None,
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct AuctionResponse {
    /// List of ad units returned in the bidding
    #[serde(rename = "ad_units")]
    pub ad_units: Vec<sdk::AdUnit>,

    /// ID of the bidding configuration
    #[serde(rename = "auction_configuration_id")]
    pub auction_configuration_id: i64,

    /// UID of the bidding configuration
    #[serde(rename = "auction_configuration_uid")]
    pub auction_configuration_uid: String,

    /// Unique identifier for the bidding
    #[serde(rename = "auction_id")]
    pub auction_id: String,

    /// PriceFloor for the bidding
    #[serde(rename = "auction_pricefloor")]
    #[validate(range(min = 0))]
    pub auction_pricefloor: f64,

    /// Timeout for the bidding in milliseconds
    #[serde(rename = "auction_timeout")]
    pub auction_timeout: i32,

    /// Indicates if external win notifications are enabled
    #[serde(rename = "external_win_notifications")]
    pub external_win_notifications: bool,

    /// List of ad units that received no bids
    #[serde(rename = "no_bids")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub no_bids: Option<Vec<sdk::AdUnit>>,

    #[serde(rename = "segment")]
    pub segment: sdk::Segment,

    /// Token
    #[serde(rename = "token")]
    pub token: String,
}

impl AuctionResponse {
    #[allow(clippy::new_without_default)]
    pub fn new(
        ad_units: Vec<sdk::AdUnit>,
        auction_configuration_id: i64,
        auction_configuration_uid: String,
        auction_id: String,
        auction_pricefloor: f64,
        auction_timeout: i32,
        external_win_notifications: bool,
        segment: sdk::Segment,
        token: String,
    ) -> AuctionResponse {
        AuctionResponse {
            ad_units,
            auction_configuration_id,
            auction_configuration_uid,
            auction_id,
            auction_pricefloor,
            auction_timeout,
            external_win_notifications,
            no_bids: None,
            segment,
            token,
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct AuctionResult {
    /// Timestamp when the bidding finished
    #[serde(rename = "auction_finish_ts")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub auction_finish_ts: Option<i64>,

    /// Timestamp when the bidding started
    #[serde(rename = "auction_start_ts")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub auction_start_ts: Option<i64>,

    #[serde(rename = "banner")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub banner: Option<sdk::BannerAdObject>,

    #[serde(rename = "bid_type")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub bid_type: Option<sdk::AuctionResultBidType>,

    /// Empty schema for interstitial ad configuration
    #[serde(rename = "interstitial")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub interstitial: Option<serde_json::Value>,

    /// Price of the winning bid
    #[serde(rename = "price")]
    #[validate(range(min = 0))]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub price: Option<f64>,

    /// Empty schema for rewarded ad configuration
    #[serde(rename = "rewarded")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub rewarded: Option<serde_json::Value>,

    #[serde(rename = "status")]
    pub status: sdk::AuctionResultStatus,

    /// Label of the winning ad unit, if applicable
    #[serde(rename = "winner_ad_unit_label")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub winner_ad_unit_label: Option<String>,

    /// UID of the winning ad unit, if applicable
    #[serde(rename = "winner_ad_unit_uid")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub winner_ad_unit_uid: Option<String>,

    /// ID of the winning demand source, if applicable
    #[serde(rename = "winner_demand_id")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub winner_demand_id: Option<String>,
}

impl AuctionResult {
    #[allow(clippy::new_without_default)]
    pub fn new(status: sdk::AuctionResultStatus) -> AuctionResult {
        AuctionResult {
            auction_finish_ts: None,
            auction_start_ts: None,
            banner: None,
            bid_type: None,
            interstitial: None,
            price: None,
            rewarded: None,
            status,
            winner_ad_unit_label: None,
            winner_ad_unit_uid: None,
            winner_demand_id: None,
        }
    }
}

/// Type of bid (RTB or CPM)
/// Enumeration of values.
/// Since this enum's variants do not hold data, we can easily define them as `#[repr(C)]`
/// which helps with FFI.
#[allow(non_camel_case_types)]
#[repr(C)]
#[derive(
    Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, serde::Serialize, serde::Deserialize, Hash,
)]
pub enum AuctionResultBidType {
    #[serde(rename = "RTB")]
    Rtb,
    #[serde(rename = "CPM")]
    Cpm,
}

/// Status of the bidding
/// Enumeration of values.
/// Since this enum's variants do not hold data, we can easily define them as `#[repr(C)]`
/// which helps with FFI.
#[allow(non_camel_case_types)]
#[repr(C)]
#[derive(
    Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, serde::Serialize, serde::Deserialize, Hash,
)]
pub enum AuctionResultStatus {
    #[serde(rename = "SUCCESS")]
    Success,
    #[serde(rename = "FAIL")]
    Fail,
    #[serde(rename = "AUCTION_CANCELLED")]
    AuctionCancelled,
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct BannerAdObject {
    #[serde(rename = "format")]
    pub format: sdk::AdFormat,
}

impl BannerAdObject {
    #[allow(clippy::new_without_default)]
    pub fn new(format: sdk::AdFormat) -> BannerAdObject {
        BannerAdObject { format }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct BaseRequest {
    #[serde(rename = "app")]
    pub app: sdk::App,

    #[serde(rename = "device")]
    pub device: sdk::Device,

    #[serde(rename = "ext")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ext: Option<String>,

    #[serde(rename = "geo")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub geo: Option<sdk::Geo>,

    #[serde(rename = "regs")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub regs: Option<sdk::Regulations>,

    #[serde(rename = "segment")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub segment: Option<sdk::Segment>,

    #[serde(rename = "session")]
    pub session: sdk::Session,

    #[serde(rename = "token")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub token: Option<String>,

    #[serde(rename = "user")]
    pub user: sdk::User,
}

impl BaseRequest {
    #[allow(clippy::new_without_default)]
    pub fn new(
        app: sdk::App,
        device: sdk::Device,
        session: sdk::Session,
        user: sdk::User,
    ) -> BaseRequest {
        BaseRequest {
            app,
            device,
            ext: None,
            geo: None,
            regs: None,
            segment: None,
            session,
            token: None,
            user,
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize)]
pub struct Bid();

/// Type of bid (RTB or CPM)
/// Enumeration of values.
/// Since this enum's variants do not hold data, we can easily define them as `#[repr(C)]`
/// which helps with FFI.
#[allow(non_camel_case_types)]
#[repr(C)]
#[derive(
    Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, serde::Serialize, serde::Deserialize, Hash,
)]
pub enum BidType {
    #[serde(rename = "RTB")]
    Rtb,
    #[serde(rename = "CPM")]
    Cpm,
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct BiddingConfig {
    /// The timeout duration for the token in milliseconds.
    #[serde(rename = "token_timeout_ms")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub token_timeout_ms: Option<i64>,
}

impl BiddingConfig {
    #[allow(clippy::new_without_default)]
    pub fn new() -> BiddingConfig {
        BiddingConfig {
            token_timeout_ms: None,
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct ClickRequest {
    #[serde(rename = "app")]
    pub app: sdk::App,

    #[serde(rename = "device")]
    pub device: sdk::Device,

    #[serde(rename = "ext")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ext: Option<String>,

    #[serde(rename = "geo")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub geo: Option<sdk::Geo>,

    #[serde(rename = "regs")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub regs: Option<sdk::Regulations>,

    #[serde(rename = "segment")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub segment: Option<sdk::Segment>,

    #[serde(rename = "session")]
    pub session: sdk::Session,

    #[serde(rename = "token")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub token: Option<String>,

    #[serde(rename = "user")]
    pub user: sdk::User,

    #[serde(rename = "bid")]
    #[serde(deserialize_with = "swagger::nullable_format::deserialize_optional_nullable")]
    #[serde(default = "swagger::nullable_format::default_optional_nullable")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub bid: Option<swagger::Nullable<sdk::Bid>>,

    #[serde(rename = "show")]
    #[serde(deserialize_with = "swagger::nullable_format::deserialize_optional_nullable")]
    #[serde(default = "swagger::nullable_format::default_optional_nullable")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub show: Option<swagger::Nullable<sdk::Bid>>,
}

impl ClickRequest {
    #[allow(clippy::new_without_default)]
    pub fn new(
        app: sdk::App,
        device: sdk::Device,
        session: sdk::Session,
        user: sdk::User,
    ) -> ClickRequest {
        ClickRequest {
            app,
            device,
            ext: None,
            geo: None,
            regs: None,
            segment: None,
            session,
            token: None,
            user,
            bid: None,
            show: None,
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct ConfigRequest {
    #[serde(rename = "app")]
    pub app: sdk::App,

    #[serde(rename = "device")]
    pub device: sdk::Device,

    #[serde(rename = "ext")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ext: Option<String>,

    #[serde(rename = "geo")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub geo: Option<sdk::Geo>,

    #[serde(rename = "regs")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub regs: Option<sdk::Regulations>,

    #[serde(rename = "segment")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub segment: Option<sdk::Segment>,

    #[serde(rename = "session")]
    pub session: sdk::Session,

    #[serde(rename = "token")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub token: Option<String>,

    #[serde(rename = "user")]
    pub user: sdk::User,

    #[serde(rename = "adapters")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub adapters: Option<std::collections::HashMap<String, sdk::Adapter>>,
}

impl ConfigRequest {
    #[allow(clippy::new_without_default)]
    pub fn new(
        app: sdk::App,
        device: sdk::Device,
        session: sdk::Session,
        user: sdk::User,
    ) -> ConfigRequest {
        ConfigRequest {
            app,
            device,
            ext: None,
            geo: None,
            regs: None,
            segment: None,
            session,
            token: None,
            user,
            adapters: None,
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct ConfigResponse {
    #[serde(rename = "bidding")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub bidding: Option<sdk::BiddingConfig>,

    #[serde(rename = "init")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub init: Option<sdk::InitConfig>,

    #[serde(rename = "placements")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub placements: Option<Vec<serde_json::Value>>,

    #[serde(rename = "segment")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub segment: Option<sdk::Segment>,

    #[serde(rename = "token")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub token: Option<String>,
}

impl ConfigResponse {
    #[allow(clippy::new_without_default)]
    pub fn new() -> ConfigResponse {
        ConfigResponse {
            bidding: None,
            init: None,
            placements: None,
            segment: None,
            token: None,
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct Device {
    /// Carrier
    #[serde(rename = "carrier")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub carrier: Option<String>,

    #[serde(rename = "connection_type")]
    pub connection_type: sdk::DeviceConnectionType,

    #[serde(rename = "geo")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub geo: Option<sdk::Geo>,

    /// Height
    #[serde(rename = "h")]
    pub h: i32,

    /// Hardware Version
    #[serde(rename = "hwv")]
    pub hwv: String,

    /// JavaScript support
    #[serde(rename = "js")]
    pub js: i32,

    /// Language
    #[serde(rename = "language")]
    pub language: String,

    /// Manufacturer
    #[serde(rename = "make")]
    pub make: String,

    /// Mobile Country Code and Mobile Network Code
    #[serde(rename = "mccmnc")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub mccmnc: Option<String>,

    /// Model
    #[serde(rename = "model")]
    pub model: String,

    /// Operating System
    #[serde(rename = "os")]
    pub os: String,

    /// Operating System Version
    #[serde(rename = "osv")]
    pub osv: String,

    /// Pixels per Inch (PPI)
    #[serde(rename = "ppi")]
    pub ppi: i32,

    /// Pixel Ratio
    #[serde(rename = "pxratio")]
    pub pxratio: f64,

    #[serde(rename = "type")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub device_type: Option<sdk::DeviceType>,

    /// UserAgent
    #[serde(rename = "ua")]
    pub ua: String,

    /// Width
    #[serde(rename = "w")]
    pub w: i32,
}

impl Device {
    #[allow(clippy::new_without_default)]
    pub fn new(
        connection_type: sdk::DeviceConnectionType,
        h: i32,
        hwv: String,
        js: i32,
        language: String,
        make: String,
        model: String,
        os: String,
        osv: String,
        ppi: i32,
        pxratio: f64,
        ua: String,
        w: i32,
    ) -> Device {
        Device {
            carrier: None,
            connection_type,
            geo: None,
            h,
            hwv,
            js,
            language,
            make,
            mccmnc: None,
            model,
            os,
            osv,
            ppi,
            pxratio,
            device_type: None,
            ua,
            w,
        }
    }
}

/// Connection Type
/// Enumeration of values.
/// Since this enum's variants do not hold data, we can easily define them as `#[repr(C)]`
/// which helps with FFI.
#[allow(non_camel_case_types)]
#[repr(C)]
#[derive(
    Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, serde::Serialize, serde::Deserialize, Hash,
)]
pub enum DeviceConnectionType {
    #[serde(rename = "ETHERNET")]
    Ethernet,
    #[serde(rename = "WIFI")]
    Wifi,
    #[serde(rename = "CELLULAR")]
    Cellular,
    #[serde(rename = "CELLULAR_UNKNOWN")]
    CellularUnknown,
    #[serde(rename = "CELLULAR_2_G")]
    Cellular2G,
    #[serde(rename = "CELLULAR_3_G")]
    Cellular3G,
    #[serde(rename = "CELLULAR_4_G")]
    Cellular4G,
    #[serde(rename = "CELLULAR_5_G")]
    Cellular5G,
}

/// Device Type
/// Enumeration of values.
/// Since this enum's variants do not hold data, we can easily define them as `#[repr(C)]`
/// which helps with FFI.
#[allow(non_camel_case_types)]
#[repr(C)]
#[derive(
    Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, serde::Serialize, serde::Deserialize, Hash,
)]
pub enum DeviceType {
    #[serde(rename = "PHONE")]
    Phone,
    #[serde(rename = "TABLET")]
    Tablet,
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct Error {
    #[serde(rename = "error")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub error: Option<sdk::ErrorError>,
}

impl Error {
    #[allow(clippy::new_without_default)]
    pub fn new() -> Error {
        Error { error: None }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct ErrorError {
    #[serde(rename = "code")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub code: Option<i32>,

    #[serde(rename = "message")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub message: Option<String>,
}

impl ErrorError {
    #[allow(clippy::new_without_default)]
    pub fn new() -> ErrorError {
        ErrorError {
            code: None,
            message: None,
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct ExternalWinner {
    /// Identifier for the demand source of the external winner
    #[serde(rename = "demand_id")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub demand_id: Option<String>,

    /// Effective cost per mille for the external winner
    #[serde(rename = "price")]
    #[validate(range(min = 0))]
    pub price: f64,
}

impl ExternalWinner {
    #[allow(clippy::new_without_default)]
    pub fn new(price: f64) -> ExternalWinner {
        ExternalWinner {
            demand_id: None,
            price,
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct Geo {
    /// Accuracy of the location data
    #[serde(rename = "accuracy")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub accuracy: Option<f64>,

    /// City name
    #[serde(rename = "city")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub city: Option<String>,

    /// Country code or name
    #[serde(rename = "country")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub country: Option<String>,

    /// Timestamp of the last location fix in seconds since epoch
    #[serde(rename = "lastfix")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub lastfix: Option<i32>,

    /// Latitude of the location
    #[serde(rename = "lat")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub lat: Option<f64>,

    /// Longitude of the location
    #[serde(rename = "lon")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub lon: Option<f64>,

    /// UTC offset in minutes
    #[serde(rename = "utcoffset")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub utcoffset: Option<i32>,

    /// ZIP or postal code
    #[serde(rename = "zip")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub zip: Option<String>,
}

impl Geo {
    #[allow(clippy::new_without_default)]
    pub fn new() -> Geo {
        Geo {
            accuracy: None,
            city: None,
            country: None,
            lastfix: None,
            lat: None,
            lon: None,
            utcoffset: None,
            zip: None,
        }
    }
}

/// Enumeration of values.
/// Since this enum's variants do not hold data, we can easily define them as `#[repr(C)]`
/// which helps with FFI.
#[allow(non_camel_case_types)]
#[repr(C)]
#[derive(
    Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, serde::Serialize, serde::Deserialize, Hash,
)]
pub enum GetAuctionAdTypeParameter {
    #[serde(rename = "banner")]
    Banner,
    #[serde(rename = "interstitial")]
    Interstitial,
    #[serde(rename = "rewarded")]
    Rewarded,
}

impl std::fmt::Display for GetAuctionAdTypeParameter {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match *self {
            GetAuctionAdTypeParameter::Banner => write!(f, "banner"),
            GetAuctionAdTypeParameter::Interstitial => write!(f, "interstitial"),
            GetAuctionAdTypeParameter::Rewarded => write!(f, "rewarded"),
        }
    }
}

impl std::str::FromStr for GetAuctionAdTypeParameter {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        match s {
            "banner" => std::result::Result::Ok(GetAuctionAdTypeParameter::Banner),
            "interstitial" => std::result::Result::Ok(GetAuctionAdTypeParameter::Interstitial),
            "rewarded" => std::result::Result::Ok(GetAuctionAdTypeParameter::Rewarded),
            _ => std::result::Result::Err(format!("Value not valid: {}", s)),
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct InitConfig {
    #[serde(rename = "adapters")]
    pub adapters: std::collections::HashMap<String, serde_json::Value>,

    #[serde(rename = "tmax")]
    pub tmax: i64,
}

impl InitConfig {
    #[allow(clippy::new_without_default)]
    pub fn new(
        adapters: std::collections::HashMap<String, serde_json::Value>,
        tmax: i64,
    ) -> InitConfig {
        InitConfig { adapters, tmax }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct LossRequest {
    #[serde(rename = "app")]
    pub app: sdk::App,

    #[serde(rename = "device")]
    pub device: sdk::Device,

    #[serde(rename = "ext")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ext: Option<String>,

    #[serde(rename = "geo")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub geo: Option<sdk::Geo>,

    #[serde(rename = "regs")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub regs: Option<sdk::Regulations>,

    #[serde(rename = "segment")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub segment: Option<sdk::Segment>,

    #[serde(rename = "session")]
    pub session: sdk::Session,

    #[serde(rename = "token")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub token: Option<String>,

    #[serde(rename = "user")]
    pub user: sdk::User,

    #[serde(rename = "bid")]
    #[serde(deserialize_with = "swagger::nullable_format::deserialize_optional_nullable")]
    #[serde(default = "swagger::nullable_format::default_optional_nullable")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub bid: Option<swagger::Nullable<sdk::Bid>>,

    #[serde(rename = "show")]
    #[serde(deserialize_with = "swagger::nullable_format::deserialize_optional_nullable")]
    #[serde(default = "swagger::nullable_format::default_optional_nullable")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub show: Option<swagger::Nullable<sdk::Bid>>,

    #[serde(rename = "external_winner")]
    pub external_winner: sdk::ExternalWinner,
}

impl LossRequest {
    #[allow(clippy::new_without_default)]
    pub fn new(
        app: sdk::App,
        device: sdk::Device,
        session: sdk::Session,
        user: sdk::User,
        external_winner: sdk::ExternalWinner,
    ) -> LossRequest {
        LossRequest {
            app,
            device,
            ext: None,
            geo: None,
            regs: None,
            segment: None,
            session,
            token: None,
            user,
            bid: None,
            show: None,
            external_winner,
        }
    }
}

/// Enumeration of values.
/// Since this enum's variants do not hold data, we can easily define them as `#[repr(C)]`
/// which helps with FFI.
#[allow(non_camel_case_types)]
#[repr(C)]
#[derive(
    Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, serde::Serialize, serde::Deserialize, Hash,
)]
pub enum PostRewardAdTypeParameter {
    #[serde(rename = "rewarded")]
    Rewarded,
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct Regulations {
    /// Indicates if COPPA regulations apply
    #[serde(rename = "coppa")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub coppa: Option<bool>,

    /// EU privacy string indicating compliance
    #[serde(rename = "eu_privacy")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub eu_privacy: Option<String>,

    /// Indicates if GDPR regulations apply
    #[serde(rename = "gdpr")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub gdpr: Option<bool>,

    /// IAB-specific settings or values
    #[serde(rename = "iab")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub iab: Option<std::collections::HashMap<String, serde_json::Value>>,

    /// US privacy string indicating compliance
    #[serde(rename = "us_privacy")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub us_privacy: Option<String>,
}

impl Regulations {
    #[allow(clippy::new_without_default)]
    pub fn new() -> Regulations {
        Regulations {
            coppa: None,
            eu_privacy: None,
            gdpr: None,
            iab: None,
            us_privacy: None,
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct RewardRequest {
    #[serde(rename = "app")]
    pub app: sdk::App,

    #[serde(rename = "device")]
    pub device: sdk::Device,

    #[serde(rename = "ext")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ext: Option<String>,

    #[serde(rename = "geo")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub geo: Option<sdk::Geo>,

    #[serde(rename = "regs")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub regs: Option<sdk::Regulations>,

    #[serde(rename = "segment")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub segment: Option<sdk::Segment>,

    #[serde(rename = "session")]
    pub session: sdk::Session,

    #[serde(rename = "token")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub token: Option<String>,

    #[serde(rename = "user")]
    pub user: sdk::User,

    #[serde(rename = "bid")]
    #[serde(deserialize_with = "swagger::nullable_format::deserialize_optional_nullable")]
    #[serde(default = "swagger::nullable_format::default_optional_nullable")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub bid: Option<swagger::Nullable<sdk::Bid>>,

    #[serde(rename = "show")]
    #[serde(deserialize_with = "swagger::nullable_format::deserialize_optional_nullable")]
    #[serde(default = "swagger::nullable_format::default_optional_nullable")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub show: Option<swagger::Nullable<sdk::Bid>>,
}

impl RewardRequest {
    #[allow(clippy::new_without_default)]
    pub fn new(
        app: sdk::App,
        device: sdk::Device,
        session: sdk::Session,
        user: sdk::User,
    ) -> RewardRequest {
        RewardRequest {
            app,
            device,
            ext: None,
            geo: None,
            regs: None,
            segment: None,
            session,
            token: None,
            user,
            bid: None,
            show: None,
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct Segment {
    /// An extension field for additional information about the segment.
    #[serde(rename = "ext")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ext: Option<String>,

    #[serde(rename = "id")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub id: Option<String>,

    #[serde(rename = "uid")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub uid: Option<String>,
}

impl Segment {
    #[allow(clippy::new_without_default)]
    pub fn new() -> Segment {
        Segment {
            ext: None,
            id: None,
            uid: None,
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct Session {
    /// Battery level percentage
    #[serde(rename = "battery")]
    pub battery: f64,

    /// CPU usage percentage
    #[serde(rename = "cpu_usage")]
    pub cpu_usage: f64,

    /// Unique identifier for the session
    #[serde(rename = "id")]
    pub id: uuid::Uuid,

    /// Monotonic timestamp of the session launch
    #[serde(rename = "launch_monotonic_ts")]
    pub launch_monotonic_ts: i64,

    /// Timestamp of the session launch
    #[serde(rename = "launch_ts")]
    pub launch_ts: i64,

    /// Monotonic timestamps when memory warnings occurred
    #[serde(rename = "memory_warnings_monotonic_ts")]
    pub memory_warnings_monotonic_ts: Vec<i64>,

    /// Timestamps when memory warnings occurred
    #[serde(rename = "memory_warnings_ts")]
    pub memory_warnings_ts: Vec<i64>,

    /// Current monotonic timestamp of the session
    #[serde(rename = "monotonic_ts")]
    pub monotonic_ts: i64,

    /// Total size of RAM
    #[serde(rename = "ram_size")]
    pub ram_size: i64,

    /// Amount of RAM used
    #[serde(rename = "ram_used")]
    pub ram_used: i64,

    /// Monotonic timestamp of the session start
    #[serde(rename = "start_monotonic_ts")]
    pub start_monotonic_ts: i64,

    /// Timestamp of the session start
    #[serde(rename = "start_ts")]
    pub start_ts: i64,

    /// Free storage space available
    #[serde(rename = "storage_free")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub storage_free: Option<i64>,

    /// Used storage space
    #[serde(rename = "storage_used")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub storage_used: Option<i64>,

    /// Current timestamp of the session
    #[serde(rename = "ts")]
    pub ts: i64,
}

impl Session {
    #[allow(clippy::new_without_default)]
    pub fn new(
        battery: f64,
        cpu_usage: f64,
        id: uuid::Uuid,
        launch_monotonic_ts: i64,
        launch_ts: i64,
        memory_warnings_monotonic_ts: Vec<i64>,
        memory_warnings_ts: Vec<i64>,
        monotonic_ts: i64,
        ram_size: i64,
        ram_used: i64,
        start_monotonic_ts: i64,
        start_ts: i64,
        storage_free: Option<i64>,
        storage_used: Option<i64>,
        ts: i64,
    ) -> Session {
        Session {
            battery,
            cpu_usage,
            id,
            launch_monotonic_ts,
            launch_ts,
            memory_warnings_monotonic_ts,
            memory_warnings_ts,
            monotonic_ts,
            ram_size,
            ram_used,
            start_monotonic_ts,
            start_ts,
            storage_free,
            storage_used,
            ts,
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct ShowRequest {
    #[serde(rename = "app")]
    pub app: sdk::App,

    #[serde(rename = "device")]
    pub device: sdk::Device,

    #[serde(rename = "ext")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ext: Option<String>,

    #[serde(rename = "geo")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub geo: Option<sdk::Geo>,

    #[serde(rename = "regs")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub regs: Option<sdk::Regulations>,

    #[serde(rename = "segment")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub segment: Option<sdk::Segment>,

    #[serde(rename = "session")]
    pub session: sdk::Session,

    #[serde(rename = "token")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub token: Option<String>,

    #[serde(rename = "user")]
    pub user: sdk::User,

    #[serde(rename = "bid")]
    #[serde(deserialize_with = "swagger::nullable_format::deserialize_optional_nullable")]
    #[serde(default = "swagger::nullable_format::default_optional_nullable")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub bid: Option<swagger::Nullable<sdk::Bid>>,

    #[serde(rename = "show")]
    #[serde(deserialize_with = "swagger::nullable_format::deserialize_optional_nullable")]
    #[serde(default = "swagger::nullable_format::default_optional_nullable")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub show: Option<swagger::Nullable<sdk::Bid>>,
}

impl ShowRequest {
    #[allow(clippy::new_without_default)]
    pub fn new(
        app: sdk::App,
        device: sdk::Device,
        session: sdk::Session,
        user: sdk::User,
    ) -> ShowRequest {
        ShowRequest {
            app,
            device,
            ext: None,
            geo: None,
            regs: None,
            segment: None,
            session,
            token: None,
            user,
            bid: None,
            show: None,
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize)]
pub struct Stats();

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct StatsRequest {
    #[serde(rename = "app")]
    pub app: sdk::App,

    #[serde(rename = "device")]
    pub device: sdk::Device,

    #[serde(rename = "ext")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ext: Option<String>,

    #[serde(rename = "geo")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub geo: Option<sdk::Geo>,

    #[serde(rename = "regs")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub regs: Option<sdk::Regulations>,

    #[serde(rename = "segment")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub segment: Option<sdk::Segment>,

    #[serde(rename = "session")]
    pub session: sdk::Session,

    #[serde(rename = "token")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub token: Option<String>,

    #[serde(rename = "user")]
    pub user: sdk::User,

    #[serde(rename = "stats")]
    pub stats: swagger::Nullable<sdk::Stats>,
}

impl StatsRequest {
    #[allow(clippy::new_without_default)]
    pub fn new(
        app: sdk::App,
        device: sdk::Device,
        session: sdk::Session,
        user: sdk::User,
        stats: swagger::Nullable<sdk::Stats>,
    ) -> StatsRequest {
        StatsRequest {
            app,
            device,
            ext: None,
            geo: None,
            regs: None,
            segment: None,
            session,
            token: None,
            user,
            stats,
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct SuccessResponse {
    /// Indicates if the operation was successful
    #[serde(rename = "success")]
    pub success: bool,
}

impl SuccessResponse {
    #[allow(clippy::new_without_default)]
    pub fn new(success: bool) -> SuccessResponse {
        SuccessResponse { success }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct User {
    /// Consent settings or preferences
    #[serde(rename = "consent")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub consent: Option<std::collections::HashMap<String, serde_json::Value>>,

    /// Indicates if COPPA (Children's Online Privacy Protection Act) applies
    #[serde(rename = "coppa")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub coppa: Option<bool>,

    /// Identifier for Advertisers (IDFA)
    #[serde(rename = "idfa")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub idfa: Option<uuid::Uuid>,

    /// Identifier for Vendors (IDFV)
    #[serde(rename = "idfv")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub idfv: Option<uuid::Uuid>,

    /// Generic identifier (IDG)
    #[serde(rename = "idg")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub idg: Option<uuid::Uuid>,

    /// Status of tracking authorization
    #[serde(rename = "tracking_authorization_status")]
    pub tracking_authorization_status: String,
}

impl User {
    #[allow(clippy::new_without_default)]
    pub fn new(tracking_authorization_status: String) -> User {
        User {
            consent: None,
            coppa: None,
            idfa: None,
            idfv: None,
            idg: None,
            tracking_authorization_status,
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct WinRequest {
    #[serde(rename = "app")]
    pub app: sdk::App,

    #[serde(rename = "device")]
    pub device: sdk::Device,

    #[serde(rename = "ext")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ext: Option<String>,

    #[serde(rename = "geo")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub geo: Option<sdk::Geo>,

    #[serde(rename = "regs")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub regs: Option<sdk::Regulations>,

    #[serde(rename = "segment")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub segment: Option<sdk::Segment>,

    #[serde(rename = "session")]
    pub session: sdk::Session,

    #[serde(rename = "token")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub token: Option<String>,

    #[serde(rename = "user")]
    pub user: sdk::User,

    #[serde(rename = "bid")]
    #[serde(deserialize_with = "swagger::nullable_format::deserialize_optional_nullable")]
    #[serde(default = "swagger::nullable_format::default_optional_nullable")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub bid: Option<swagger::Nullable<sdk::Bid>>,

    #[serde(rename = "show")]
    #[serde(deserialize_with = "swagger::nullable_format::deserialize_optional_nullable")]
    #[serde(default = "swagger::nullable_format::default_optional_nullable")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub show: Option<swagger::Nullable<sdk::Bid>>,
}

impl WinRequest {
    #[allow(clippy::new_without_default)]
    pub fn new(
        app: sdk::App,
        device: sdk::Device,
        session: sdk::Session,
        user: sdk::User,
    ) -> WinRequest {
        WinRequest {
            app,
            device,
            ext: None,
            geo: None,
            regs: None,
            segment: None,
            session,
            token: None,
            user,
            bid: None,
            show: None,
        }
    }
}
