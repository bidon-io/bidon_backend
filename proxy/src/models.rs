#![allow(unused_qualifications)]

use crate::header;
use crate::models;

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

impl std::fmt::Display for AdFormat {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match *self {
            AdFormat::Banner => write!(f, "BANNER"),
            AdFormat::Leaderboard => write!(f, "LEADERBOARD"),
            AdFormat::Mrec => write!(f, "MREC"),
            AdFormat::Adaptive => write!(f, "ADAPTIVE"),
        }
    }
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

// Methods for converting between header::IntoHeaderValue<AdFormat> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<AdFormat>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<AdFormat>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for AdFormat - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<AdFormat> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <AdFormat as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into AdFormat - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<AdFormat>>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<AdFormat>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<Vec<AdFormat>> {
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<AdFormat> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <AdFormat as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into AdFormat - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
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
    pub banner: Option<models::BannerAdObject>,

    /// Map of demands
    #[serde(rename = "demands")]
    pub demands: std::collections::HashMap<String, serde_json::Value>,

    /// Empty schema for interstitial ad configuration
    #[serde(rename = "interstitial")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub interstitial: Option<serde_json::Value>,

    #[serde(rename = "orientation")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub orientation: Option<models::AdObjectOrientation>,

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

/// Converts the AdObject value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for AdObject {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            self.auction_configuration_id
                .as_ref()
                .map(|auction_configuration_id| {
                    [
                        "auction_configuration_id".to_string(),
                        auction_configuration_id.to_string(),
                    ]
                    .join(",")
                }),
            self.auction_configuration_uid
                .as_ref()
                .map(|auction_configuration_uid| {
                    [
                        "auction_configuration_uid".to_string(),
                        auction_configuration_uid.to_string(),
                    ]
                    .join(",")
                }),
            self.auction_id
                .as_ref()
                .map(|auction_id| ["auction_id".to_string(), auction_id.to_string()].join(",")),
            self.auction_key
                .as_ref()
                .map(|auction_key| ["auction_key".to_string(), auction_key.to_string()].join(",")),
            Some("auction_pricefloor".to_string()),
            Some(self.auction_pricefloor.to_string()),
            // Skipping non-primitive type banner in query parameter serialization
            // Skipping map demands in query parameter serialization
            // Skipping non-primitive type interstitial in query parameter serialization
            // Skipping non-primitive type orientation in query parameter serialization
            // Skipping non-primitive type rewarded in query parameter serialization
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a AdObject value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for AdObject {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub auction_configuration_id: Vec<i64>,
            pub auction_configuration_uid: Vec<String>,
            pub auction_id: Vec<String>,
            pub auction_key: Vec<String>,
            pub auction_pricefloor: Vec<f64>,
            pub banner: Vec<models::BannerAdObject>,
            pub demands: Vec<std::collections::HashMap<String, serde_json::Value>>,
            pub interstitial: Vec<serde_json::Value>,
            pub orientation: Vec<models::AdObjectOrientation>,
            pub rewarded: Vec<serde_json::Value>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err(
                        "Missing value while parsing AdObject".to_string(),
                    )
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    #[allow(clippy::redundant_clone)]
                    "auction_configuration_id" => intermediate_rep.auction_configuration_id.push(
                        <i64 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "auction_configuration_uid" => intermediate_rep.auction_configuration_uid.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "auction_id" => intermediate_rep.auction_id.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "auction_key" => intermediate_rep.auction_key.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "auction_pricefloor" => intermediate_rep.auction_pricefloor.push(
                        <f64 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "banner" => intermediate_rep.banner.push(
                        <models::BannerAdObject as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    "demands" => {
                        return std::result::Result::Err(
                            "Parsing a container in this style is not supported in AdObject"
                                .to_string(),
                        )
                    }
                    #[allow(clippy::redundant_clone)]
                    "interstitial" => intermediate_rep.interstitial.push(
                        <serde_json::Value as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "orientation" => intermediate_rep.orientation.push(
                        <models::AdObjectOrientation as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "rewarded" => intermediate_rep.rewarded.push(
                        <serde_json::Value as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing AdObject".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(AdObject {
            auction_configuration_id: intermediate_rep.auction_configuration_id.into_iter().next(),
            auction_configuration_uid: intermediate_rep
                .auction_configuration_uid
                .into_iter()
                .next(),
            auction_id: intermediate_rep.auction_id.into_iter().next(),
            auction_key: intermediate_rep.auction_key.into_iter().next(),
            auction_pricefloor: intermediate_rep
                .auction_pricefloor
                .into_iter()
                .next()
                .ok_or_else(|| "auction_pricefloor missing in AdObject".to_string())?,
            banner: intermediate_rep.banner.into_iter().next(),
            demands: intermediate_rep
                .demands
                .into_iter()
                .next()
                .ok_or_else(|| "demands missing in AdObject".to_string())?,
            interstitial: intermediate_rep.interstitial.into_iter().next(),
            orientation: intermediate_rep.orientation.into_iter().next(),
            rewarded: intermediate_rep.rewarded.into_iter().next(),
        })
    }
}

// Methods for converting between header::IntoHeaderValue<AdObject> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<AdObject>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<AdObject>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for AdObject - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<AdObject> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <AdObject as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into AdObject - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<AdObject>>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<AdObject>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<Vec<AdObject>> {
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<AdObject> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <AdObject as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into AdObject - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
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

impl std::fmt::Display for AdObjectOrientation {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match *self {
            AdObjectOrientation::Portrait => write!(f, "PORTRAIT"),
            AdObjectOrientation::Landscape => write!(f, "LANDSCAPE"),
        }
    }
}

impl std::str::FromStr for AdObjectOrientation {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        match s {
            "PORTRAIT" => std::result::Result::Ok(AdObjectOrientation::Portrait),
            "LANDSCAPE" => std::result::Result::Ok(AdObjectOrientation::Landscape),
            _ => std::result::Result::Err(format!("Value not valid: {}", s)),
        }
    }
}

// Methods for converting between header::IntoHeaderValue<AdObjectOrientation> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<AdObjectOrientation>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<AdObjectOrientation>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for AdObjectOrientation - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<AdObjectOrientation>
{
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <AdObjectOrientation as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into AdObjectOrientation - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<AdObjectOrientation>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<AdObjectOrientation>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<AdObjectOrientation>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values : std::vec::Vec<AdObjectOrientation> = hdr_values
                .split(',')
                .filter_map(|hdr_value| match hdr_value.trim() {
                    "" => std::option::Option::None,
                    hdr_value => std::option::Option::Some({
                        match <AdObjectOrientation as std::str::FromStr>::from_str(hdr_value) {
                            std::result::Result::Ok(value) => std::result::Result::Ok(value),
                            std::result::Result::Err(err) => std::result::Result::Err(
                                format!("Unable to convert header value '{}' into AdObjectOrientation - {}",
                                    hdr_value, err))
                        }
                    })
                }).collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
        }
    }
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

/// Converts the AdUnit value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for AdUnit {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            Some("bid_type".to_string()),
            Some(self.bid_type.to_string()),
            Some("demand_id".to_string()),
            Some(self.demand_id.to_string()),
            // Skipping map ext in query parameter serialization
            Some("label".to_string()),
            Some(self.label.to_string()),
            self.pricefloor
                .as_ref()
                .map(|pricefloor| ["pricefloor".to_string(), pricefloor.to_string()].join(",")),
            Some("uid".to_string()),
            Some(self.uid.to_string()),
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a AdUnit value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for AdUnit {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub bid_type: Vec<String>,
            pub demand_id: Vec<String>,
            pub ext: Vec<std::collections::HashMap<String, serde_json::Value>>,
            pub label: Vec<String>,
            pub pricefloor: Vec<f64>,
            pub uid: Vec<String>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err(
                        "Missing value while parsing AdUnit".to_string(),
                    )
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    #[allow(clippy::redundant_clone)]
                    "bid_type" => intermediate_rep.bid_type.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "demand_id" => intermediate_rep.demand_id.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    "ext" => {
                        return std::result::Result::Err(
                            "Parsing a container in this style is not supported in AdUnit"
                                .to_string(),
                        )
                    }
                    #[allow(clippy::redundant_clone)]
                    "label" => intermediate_rep.label.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "pricefloor" => intermediate_rep.pricefloor.push(
                        <f64 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "uid" => intermediate_rep.uid.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing AdUnit".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(AdUnit {
            bid_type: intermediate_rep
                .bid_type
                .into_iter()
                .next()
                .ok_or_else(|| "bid_type missing in AdUnit".to_string())?,
            demand_id: intermediate_rep
                .demand_id
                .into_iter()
                .next()
                .ok_or_else(|| "demand_id missing in AdUnit".to_string())?,
            ext: intermediate_rep.ext.into_iter().next(),
            label: intermediate_rep
                .label
                .into_iter()
                .next()
                .ok_or_else(|| "label missing in AdUnit".to_string())?,
            pricefloor: intermediate_rep.pricefloor.into_iter().next(),
            uid: intermediate_rep
                .uid
                .into_iter()
                .next()
                .ok_or_else(|| "uid missing in AdUnit".to_string())?,
        })
    }
}

// Methods for converting between header::IntoHeaderValue<AdUnit> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<AdUnit>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<AdUnit>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for AdUnit - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<AdUnit> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <AdUnit as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into AdUnit - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<AdUnit>>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<AdUnit>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<Vec<AdUnit>> {
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<AdUnit> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <AdUnit as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into AdUnit - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
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

/// Converts the Adapter value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for Adapter {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            Some("sdk_version".to_string()),
            Some(self.sdk_version.to_string()),
            Some("version".to_string()),
            Some(self.version.to_string()),
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a Adapter value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for Adapter {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub sdk_version: Vec<String>,
            pub version: Vec<String>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err(
                        "Missing value while parsing Adapter".to_string(),
                    )
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    #[allow(clippy::redundant_clone)]
                    "sdk_version" => intermediate_rep.sdk_version.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "version" => intermediate_rep.version.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing Adapter".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(Adapter {
            sdk_version: intermediate_rep
                .sdk_version
                .into_iter()
                .next()
                .ok_or_else(|| "sdk_version missing in Adapter".to_string())?,
            version: intermediate_rep
                .version
                .into_iter()
                .next()
                .ok_or_else(|| "version missing in Adapter".to_string())?,
        })
    }
}

// Methods for converting between header::IntoHeaderValue<Adapter> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<Adapter>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<Adapter>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for Adapter - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<Adapter> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <Adapter as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into Adapter - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<Adapter>>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<Adapter>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<Vec<Adapter>> {
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<Adapter> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <Adapter as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into Adapter - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
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

/// Converts the App value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for App {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            Some("bundle".to_string()),
            Some(self.bundle.to_string()),
            Some("framework".to_string()),
            Some(self.framework.to_string()),
            self.framework_version.as_ref().map(|framework_version| {
                [
                    "framework_version".to_string(),
                    framework_version.to_string(),
                ]
                .join(",")
            }),
            Some("key".to_string()),
            Some(self.key.to_string()),
            self.plugin_version.as_ref().map(|plugin_version| {
                ["plugin_version".to_string(), plugin_version.to_string()].join(",")
            }),
            self.sdk_version
                .as_ref()
                .map(|sdk_version| ["sdk_version".to_string(), sdk_version.to_string()].join(",")),
            self.skadn.as_ref().map(|skadn| {
                [
                    "skadn".to_string(),
                    skadn
                        .iter()
                        .map(|x| x.to_string())
                        .collect::<Vec<_>>()
                        .join(","),
                ]
                .join(",")
            }),
            Some("version".to_string()),
            Some(self.version.to_string()),
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a App value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for App {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub bundle: Vec<String>,
            pub framework: Vec<String>,
            pub framework_version: Vec<String>,
            pub key: Vec<String>,
            pub plugin_version: Vec<String>,
            pub sdk_version: Vec<String>,
            pub skadn: Vec<Vec<String>>,
            pub version: Vec<String>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err("Missing value while parsing App".to_string())
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    #[allow(clippy::redundant_clone)]
                    "bundle" => intermediate_rep.bundle.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "framework" => intermediate_rep.framework.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "framework_version" => intermediate_rep.framework_version.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "key" => intermediate_rep.key.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "plugin_version" => intermediate_rep.plugin_version.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "sdk_version" => intermediate_rep.sdk_version.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    "skadn" => {
                        return std::result::Result::Err(
                            "Parsing a container in this style is not supported in App".to_string(),
                        )
                    }
                    #[allow(clippy::redundant_clone)]
                    "version" => intermediate_rep.version.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing App".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(App {
            bundle: intermediate_rep
                .bundle
                .into_iter()
                .next()
                .ok_or_else(|| "bundle missing in App".to_string())?,
            framework: intermediate_rep
                .framework
                .into_iter()
                .next()
                .ok_or_else(|| "framework missing in App".to_string())?,
            framework_version: intermediate_rep.framework_version.into_iter().next(),
            key: intermediate_rep
                .key
                .into_iter()
                .next()
                .ok_or_else(|| "key missing in App".to_string())?,
            plugin_version: intermediate_rep.plugin_version.into_iter().next(),
            sdk_version: intermediate_rep.sdk_version.into_iter().next(),
            skadn: intermediate_rep.skadn.into_iter().next(),
            version: intermediate_rep
                .version
                .into_iter()
                .next()
                .ok_or_else(|| "version missing in App".to_string())?,
        })
    }
}

// Methods for converting between header::IntoHeaderValue<App> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<App>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(hdr_value: header::IntoHeaderValue<App>) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for App - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<App> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => match <App as std::str::FromStr>::from_str(value) {
                std::result::Result::Ok(value) => {
                    std::result::Result::Ok(header::IntoHeaderValue(value))
                }
                std::result::Result::Err(err) => std::result::Result::Err(format!(
                    "Unable to convert header value '{}' into App - {}",
                    value, err
                )),
            },
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<App>>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<App>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<Vec<App>> {
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<App> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <App as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into App - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
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
    pub status: models::AuctionAdUnitResultStatus,

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
    pub fn new(
        demand_id: String,
        status: models::AuctionAdUnitResultStatus,
    ) -> AuctionAdUnitResult {
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

/// Converts the AuctionAdUnitResult value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for AuctionAdUnitResult {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            self.ad_unit_label.as_ref().map(|ad_unit_label| {
                ["ad_unit_label".to_string(), ad_unit_label.to_string()].join(",")
            }),
            self.ad_unit_uid
                .as_ref()
                .map(|ad_unit_uid| ["ad_unit_uid".to_string(), ad_unit_uid.to_string()].join(",")),
            // Skipping non-primitive type bid_type in query parameter serialization
            Some("demand_id".to_string()),
            Some(self.demand_id.to_string()),
            self.error_message.as_ref().map(|error_message| {
                ["error_message".to_string(), error_message.to_string()].join(",")
            }),
            self.fill_finish_ts.as_ref().map(|fill_finish_ts| {
                ["fill_finish_ts".to_string(), fill_finish_ts.to_string()].join(",")
            }),
            self.fill_start_ts.as_ref().map(|fill_start_ts| {
                ["fill_start_ts".to_string(), fill_start_ts.to_string()].join(",")
            }),
            self.price
                .as_ref()
                .map(|price| ["price".to_string(), price.to_string()].join(",")),
            // Skipping non-primitive type status in query parameter serialization
            self.token_finish_ts.as_ref().map(|token_finish_ts| {
                ["token_finish_ts".to_string(), token_finish_ts.to_string()].join(",")
            }),
            self.token_start_ts.as_ref().map(|token_start_ts| {
                ["token_start_ts".to_string(), token_start_ts.to_string()].join(",")
            }),
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a AuctionAdUnitResult value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for AuctionAdUnitResult {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub ad_unit_label: Vec<String>,
            pub ad_unit_uid: Vec<String>,
            pub bid_type: Vec<serde_json::Value>,
            pub demand_id: Vec<String>,
            pub error_message: Vec<String>,
            pub fill_finish_ts: Vec<i64>,
            pub fill_start_ts: Vec<i64>,
            pub price: Vec<f64>,
            pub status: Vec<models::AuctionAdUnitResultStatus>,
            pub token_finish_ts: Vec<i64>,
            pub token_start_ts: Vec<i64>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err(
                        "Missing value while parsing AuctionAdUnitResult".to_string(),
                    )
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    #[allow(clippy::redundant_clone)]
                    "ad_unit_label" => intermediate_rep.ad_unit_label.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "ad_unit_uid" => intermediate_rep.ad_unit_uid.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "bid_type" => intermediate_rep.bid_type.push(
                        <serde_json::Value as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "demand_id" => intermediate_rep.demand_id.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "error_message" => intermediate_rep.error_message.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "fill_finish_ts" => intermediate_rep.fill_finish_ts.push(
                        <i64 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "fill_start_ts" => intermediate_rep.fill_start_ts.push(
                        <i64 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "price" => intermediate_rep.price.push(
                        <f64 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "status" => intermediate_rep.status.push(
                        <models::AuctionAdUnitResultStatus as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "token_finish_ts" => intermediate_rep.token_finish_ts.push(
                        <i64 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "token_start_ts" => intermediate_rep.token_start_ts.push(
                        <i64 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing AuctionAdUnitResult".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(AuctionAdUnitResult {
            ad_unit_label: intermediate_rep.ad_unit_label.into_iter().next(),
            ad_unit_uid: intermediate_rep.ad_unit_uid.into_iter().next(),
            bid_type: intermediate_rep.bid_type.into_iter().next(),
            demand_id: intermediate_rep
                .demand_id
                .into_iter()
                .next()
                .ok_or_else(|| "demand_id missing in AuctionAdUnitResult".to_string())?,
            error_message: intermediate_rep.error_message.into_iter().next(),
            fill_finish_ts: intermediate_rep.fill_finish_ts.into_iter().next(),
            fill_start_ts: intermediate_rep.fill_start_ts.into_iter().next(),
            price: intermediate_rep.price.into_iter().next(),
            status: intermediate_rep
                .status
                .into_iter()
                .next()
                .ok_or_else(|| "status missing in AuctionAdUnitResult".to_string())?,
            token_finish_ts: intermediate_rep.token_finish_ts.into_iter().next(),
            token_start_ts: intermediate_rep.token_start_ts.into_iter().next(),
        })
    }
}

// Methods for converting between header::IntoHeaderValue<AuctionAdUnitResult> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<AuctionAdUnitResult>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<AuctionAdUnitResult>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for AuctionAdUnitResult - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<AuctionAdUnitResult>
{
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <AuctionAdUnitResult as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into AuctionAdUnitResult - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<AuctionAdUnitResult>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<AuctionAdUnitResult>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<AuctionAdUnitResult>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values : std::vec::Vec<AuctionAdUnitResult> = hdr_values
                .split(',')
                .filter_map(|hdr_value| match hdr_value.trim() {
                    "" => std::option::Option::None,
                    hdr_value => std::option::Option::Some({
                        match <AuctionAdUnitResult as std::str::FromStr>::from_str(hdr_value) {
                            std::result::Result::Ok(value) => std::result::Result::Ok(value),
                            std::result::Result::Err(err) => std::result::Result::Err(
                                format!("Unable to convert header value '{}' into AuctionAdUnitResult - {}",
                                    hdr_value, err))
                        }
                    })
                }).collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
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

impl std::fmt::Display for AuctionAdUnitResultStatus {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match *self {
            AuctionAdUnitResultStatus::Empty => write!(f, ""),
            AuctionAdUnitResultStatus::Success => write!(f, "SUCCESS"),
            AuctionAdUnitResultStatus::Fail => write!(f, "FAIL"),
            AuctionAdUnitResultStatus::Pending => write!(f, "PENDING"),
        }
    }
}

impl std::str::FromStr for AuctionAdUnitResultStatus {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        match s {
            "" => std::result::Result::Ok(AuctionAdUnitResultStatus::Empty),
            "SUCCESS" => std::result::Result::Ok(AuctionAdUnitResultStatus::Success),
            "FAIL" => std::result::Result::Ok(AuctionAdUnitResultStatus::Fail),
            "PENDING" => std::result::Result::Ok(AuctionAdUnitResultStatus::Pending),
            _ => std::result::Result::Err(format!("Value not valid: {}", s)),
        }
    }
}

// Methods for converting between header::IntoHeaderValue<AuctionAdUnitResultStatus> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<AuctionAdUnitResultStatus>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<AuctionAdUnitResultStatus>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for AuctionAdUnitResultStatus - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<AuctionAdUnitResultStatus>
{
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <AuctionAdUnitResultStatus as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into AuctionAdUnitResultStatus - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<AuctionAdUnitResultStatus>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<AuctionAdUnitResultStatus>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<AuctionAdUnitResultStatus>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values : std::vec::Vec<AuctionAdUnitResultStatus> = hdr_values
                .split(',')
                .filter_map(|hdr_value| match hdr_value.trim() {
                    "" => std::option::Option::None,
                    hdr_value => std::option::Option::Some({
                        match <AuctionAdUnitResultStatus as std::str::FromStr>::from_str(hdr_value) {
                            std::result::Result::Ok(value) => std::result::Result::Ok(value),
                            std::result::Result::Err(err) => std::result::Result::Err(
                                format!("Unable to convert header value '{}' into AuctionAdUnitResultStatus - {}",
                                    hdr_value, err))
                        }
                    })
                }).collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct AuctionRequest {
    #[serde(rename = "app")]
    pub app: models::App,

    #[serde(rename = "device")]
    pub device: models::Device,

    #[serde(rename = "ext")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ext: Option<String>,

    #[serde(rename = "geo")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub geo: Option<models::Geo>,

    #[serde(rename = "regs")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub regs: Option<models::Regulations>,

    #[serde(rename = "segment")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub segment: Option<models::Segment>,

    #[serde(rename = "session")]
    pub session: models::Session,

    #[serde(rename = "token")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub token: Option<String>,

    #[serde(rename = "user")]
    pub user: models::User,

    #[serde(rename = "ad_object")]
    pub ad_object: models::AdObject,

    #[serde(rename = "adapters")]
    pub adapters: std::collections::HashMap<String, models::Adapter>,

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
        app: models::App,
        device: models::Device,
        session: models::Session,
        user: models::User,
        ad_object: models::AdObject,
        adapters: std::collections::HashMap<String, models::Adapter>,
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

/// Converts the AuctionRequest value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for AuctionRequest {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            // Skipping non-primitive type app in query parameter serialization
            // Skipping non-primitive type device in query parameter serialization
            self.ext
                .as_ref()
                .map(|ext| ["ext".to_string(), ext.to_string()].join(",")),
            // Skipping non-primitive type geo in query parameter serialization
            // Skipping non-primitive type regs in query parameter serialization
            // Skipping non-primitive type segment in query parameter serialization
            // Skipping non-primitive type session in query parameter serialization
            self.token
                .as_ref()
                .map(|token| ["token".to_string(), token.to_string()].join(",")),
            // Skipping non-primitive type user in query parameter serialization
            // Skipping non-primitive type ad_object in query parameter serialization
            // Skipping map adapters in query parameter serialization
            self.test
                .as_ref()
                .map(|test| ["test".to_string(), test.to_string()].join(",")),
            self.tmax
                .as_ref()
                .map(|tmax| ["tmax".to_string(), tmax.to_string()].join(",")),
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a AuctionRequest value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for AuctionRequest {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub app: Vec<models::App>,
            pub device: Vec<models::Device>,
            pub ext: Vec<String>,
            pub geo: Vec<models::Geo>,
            pub regs: Vec<models::Regulations>,
            pub segment: Vec<models::Segment>,
            pub session: Vec<models::Session>,
            pub token: Vec<String>,
            pub user: Vec<models::User>,
            pub ad_object: Vec<models::AdObject>,
            pub adapters: Vec<std::collections::HashMap<String, models::Adapter>>,
            pub test: Vec<bool>,
            pub tmax: Vec<i64>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err(
                        "Missing value while parsing AuctionRequest".to_string(),
                    )
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    #[allow(clippy::redundant_clone)]
                    "app" => intermediate_rep.app.push(
                        <models::App as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "device" => intermediate_rep.device.push(
                        <models::Device as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "ext" => intermediate_rep.ext.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "geo" => intermediate_rep.geo.push(
                        <models::Geo as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "regs" => intermediate_rep.regs.push(
                        <models::Regulations as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "segment" => intermediate_rep.segment.push(
                        <models::Segment as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "session" => intermediate_rep.session.push(
                        <models::Session as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "token" => intermediate_rep.token.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "user" => intermediate_rep.user.push(
                        <models::User as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "ad_object" => intermediate_rep.ad_object.push(
                        <models::AdObject as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    "adapters" => {
                        return std::result::Result::Err(
                            "Parsing a container in this style is not supported in AuctionRequest"
                                .to_string(),
                        )
                    }
                    #[allow(clippy::redundant_clone)]
                    "test" => intermediate_rep.test.push(
                        <bool as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "tmax" => intermediate_rep.tmax.push(
                        <i64 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing AuctionRequest".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(AuctionRequest {
            app: intermediate_rep
                .app
                .into_iter()
                .next()
                .ok_or_else(|| "app missing in AuctionRequest".to_string())?,
            device: intermediate_rep
                .device
                .into_iter()
                .next()
                .ok_or_else(|| "device missing in AuctionRequest".to_string())?,
            ext: intermediate_rep.ext.into_iter().next(),
            geo: intermediate_rep.geo.into_iter().next(),
            regs: intermediate_rep.regs.into_iter().next(),
            segment: intermediate_rep.segment.into_iter().next(),
            session: intermediate_rep
                .session
                .into_iter()
                .next()
                .ok_or_else(|| "session missing in AuctionRequest".to_string())?,
            token: intermediate_rep.token.into_iter().next(),
            user: intermediate_rep
                .user
                .into_iter()
                .next()
                .ok_or_else(|| "user missing in AuctionRequest".to_string())?,
            ad_object: intermediate_rep
                .ad_object
                .into_iter()
                .next()
                .ok_or_else(|| "ad_object missing in AuctionRequest".to_string())?,
            adapters: intermediate_rep
                .adapters
                .into_iter()
                .next()
                .ok_or_else(|| "adapters missing in AuctionRequest".to_string())?,
            test: intermediate_rep.test.into_iter().next(),
            tmax: intermediate_rep.tmax.into_iter().next(),
        })
    }
}

// Methods for converting between header::IntoHeaderValue<AuctionRequest> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<AuctionRequest>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<AuctionRequest>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for AuctionRequest - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<AuctionRequest> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <AuctionRequest as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into AuctionRequest - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<AuctionRequest>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<AuctionRequest>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<AuctionRequest>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<AuctionRequest> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <AuctionRequest as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into AuctionRequest - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct AuctionResponse {
    /// List of ad units returned in the bidding
    #[serde(rename = "ad_units")]
    pub ad_units: Vec<models::AdUnit>,

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
    pub no_bids: Option<Vec<models::AdUnit>>,

    #[serde(rename = "segment")]
    pub segment: models::Segment,

    /// Token
    #[serde(rename = "token")]
    pub token: String,
}

impl AuctionResponse {
    #[allow(clippy::new_without_default)]
    pub fn new(
        ad_units: Vec<models::AdUnit>,
        auction_configuration_id: i64,
        auction_configuration_uid: String,
        auction_id: String,
        auction_pricefloor: f64,
        auction_timeout: i32,
        external_win_notifications: bool,
        segment: models::Segment,
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

/// Converts the AuctionResponse value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for AuctionResponse {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            // Skipping non-primitive type ad_units in query parameter serialization
            Some("auction_configuration_id".to_string()),
            Some(self.auction_configuration_id.to_string()),
            Some("auction_configuration_uid".to_string()),
            Some(self.auction_configuration_uid.to_string()),
            Some("auction_id".to_string()),
            Some(self.auction_id.to_string()),
            Some("auction_pricefloor".to_string()),
            Some(self.auction_pricefloor.to_string()),
            Some("auction_timeout".to_string()),
            Some(self.auction_timeout.to_string()),
            Some("external_win_notifications".to_string()),
            Some(self.external_win_notifications.to_string()),
            // Skipping non-primitive type no_bids in query parameter serialization
            // Skipping non-primitive type segment in query parameter serialization
            Some("token".to_string()),
            Some(self.token.to_string()),
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a AuctionResponse value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for AuctionResponse {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub ad_units: Vec<Vec<models::AdUnit>>,
            pub auction_configuration_id: Vec<i64>,
            pub auction_configuration_uid: Vec<String>,
            pub auction_id: Vec<String>,
            pub auction_pricefloor: Vec<f64>,
            pub auction_timeout: Vec<i32>,
            pub external_win_notifications: Vec<bool>,
            pub no_bids: Vec<Vec<models::AdUnit>>,
            pub segment: Vec<models::Segment>,
            pub token: Vec<String>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err(
                        "Missing value while parsing AuctionResponse".to_string(),
                    )
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    "ad_units" => {
                        return std::result::Result::Err(
                            "Parsing a container in this style is not supported in AuctionResponse"
                                .to_string(),
                        )
                    }
                    #[allow(clippy::redundant_clone)]
                    "auction_configuration_id" => intermediate_rep.auction_configuration_id.push(
                        <i64 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "auction_configuration_uid" => intermediate_rep.auction_configuration_uid.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "auction_id" => intermediate_rep.auction_id.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "auction_pricefloor" => intermediate_rep.auction_pricefloor.push(
                        <f64 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "auction_timeout" => intermediate_rep.auction_timeout.push(
                        <i32 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "external_win_notifications" => {
                        intermediate_rep.external_win_notifications.push(
                            <bool as std::str::FromStr>::from_str(val)
                                .map_err(|x| x.to_string())?,
                        )
                    }
                    "no_bids" => {
                        return std::result::Result::Err(
                            "Parsing a container in this style is not supported in AuctionResponse"
                                .to_string(),
                        )
                    }
                    #[allow(clippy::redundant_clone)]
                    "segment" => intermediate_rep.segment.push(
                        <models::Segment as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "token" => intermediate_rep.token.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing AuctionResponse".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(AuctionResponse {
            ad_units: intermediate_rep
                .ad_units
                .into_iter()
                .next()
                .ok_or_else(|| "ad_units missing in AuctionResponse".to_string())?,
            auction_configuration_id: intermediate_rep
                .auction_configuration_id
                .into_iter()
                .next()
                .ok_or_else(|| "auction_configuration_id missing in AuctionResponse".to_string())?,
            auction_configuration_uid: intermediate_rep
                .auction_configuration_uid
                .into_iter()
                .next()
                .ok_or_else(|| {
                    "auction_configuration_uid missing in AuctionResponse".to_string()
                })?,
            auction_id: intermediate_rep
                .auction_id
                .into_iter()
                .next()
                .ok_or_else(|| "auction_id missing in AuctionResponse".to_string())?,
            auction_pricefloor: intermediate_rep
                .auction_pricefloor
                .into_iter()
                .next()
                .ok_or_else(|| "auction_pricefloor missing in AuctionResponse".to_string())?,
            auction_timeout: intermediate_rep
                .auction_timeout
                .into_iter()
                .next()
                .ok_or_else(|| "auction_timeout missing in AuctionResponse".to_string())?,
            external_win_notifications: intermediate_rep
                .external_win_notifications
                .into_iter()
                .next()
                .ok_or_else(|| {
                    "external_win_notifications missing in AuctionResponse".to_string()
                })?,
            no_bids: intermediate_rep.no_bids.into_iter().next(),
            segment: intermediate_rep
                .segment
                .into_iter()
                .next()
                .ok_or_else(|| "segment missing in AuctionResponse".to_string())?,
            token: intermediate_rep
                .token
                .into_iter()
                .next()
                .ok_or_else(|| "token missing in AuctionResponse".to_string())?,
        })
    }
}

// Methods for converting between header::IntoHeaderValue<AuctionResponse> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<AuctionResponse>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<AuctionResponse>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for AuctionResponse - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<AuctionResponse>
{
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <AuctionResponse as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into AuctionResponse - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<AuctionResponse>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<AuctionResponse>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<AuctionResponse>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<AuctionResponse> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <AuctionResponse as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into AuctionResponse - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
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
    pub banner: Option<models::BannerAdObject>,

    #[serde(rename = "bid_type")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub bid_type: Option<models::AuctionResultBidType>,

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
    pub status: models::AuctionResultStatus,

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
    pub fn new(status: models::AuctionResultStatus) -> AuctionResult {
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

/// Converts the AuctionResult value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for AuctionResult {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            self.auction_finish_ts.as_ref().map(|auction_finish_ts| {
                [
                    "auction_finish_ts".to_string(),
                    auction_finish_ts.to_string(),
                ]
                .join(",")
            }),
            self.auction_start_ts.as_ref().map(|auction_start_ts| {
                ["auction_start_ts".to_string(), auction_start_ts.to_string()].join(",")
            }),
            // Skipping non-primitive type banner in query parameter serialization
            // Skipping non-primitive type bid_type in query parameter serialization
            // Skipping non-primitive type interstitial in query parameter serialization
            self.price
                .as_ref()
                .map(|price| ["price".to_string(), price.to_string()].join(",")),
            // Skipping non-primitive type rewarded in query parameter serialization
            // Skipping non-primitive type status in query parameter serialization
            self.winner_ad_unit_label
                .as_ref()
                .map(|winner_ad_unit_label| {
                    [
                        "winner_ad_unit_label".to_string(),
                        winner_ad_unit_label.to_string(),
                    ]
                    .join(",")
                }),
            self.winner_ad_unit_uid.as_ref().map(|winner_ad_unit_uid| {
                [
                    "winner_ad_unit_uid".to_string(),
                    winner_ad_unit_uid.to_string(),
                ]
                .join(",")
            }),
            self.winner_demand_id.as_ref().map(|winner_demand_id| {
                ["winner_demand_id".to_string(), winner_demand_id.to_string()].join(",")
            }),
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a AuctionResult value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for AuctionResult {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub auction_finish_ts: Vec<i64>,
            pub auction_start_ts: Vec<i64>,
            pub banner: Vec<models::BannerAdObject>,
            pub bid_type: Vec<models::AuctionResultBidType>,
            pub interstitial: Vec<serde_json::Value>,
            pub price: Vec<f64>,
            pub rewarded: Vec<serde_json::Value>,
            pub status: Vec<models::AuctionResultStatus>,
            pub winner_ad_unit_label: Vec<String>,
            pub winner_ad_unit_uid: Vec<String>,
            pub winner_demand_id: Vec<String>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err(
                        "Missing value while parsing AuctionResult".to_string(),
                    )
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    #[allow(clippy::redundant_clone)]
                    "auction_finish_ts" => intermediate_rep.auction_finish_ts.push(
                        <i64 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "auction_start_ts" => intermediate_rep.auction_start_ts.push(
                        <i64 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "banner" => intermediate_rep.banner.push(
                        <models::BannerAdObject as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "bid_type" => intermediate_rep.bid_type.push(
                        <models::AuctionResultBidType as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "interstitial" => intermediate_rep.interstitial.push(
                        <serde_json::Value as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "price" => intermediate_rep.price.push(
                        <f64 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "rewarded" => intermediate_rep.rewarded.push(
                        <serde_json::Value as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "status" => intermediate_rep.status.push(
                        <models::AuctionResultStatus as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "winner_ad_unit_label" => intermediate_rep.winner_ad_unit_label.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "winner_ad_unit_uid" => intermediate_rep.winner_ad_unit_uid.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "winner_demand_id" => intermediate_rep.winner_demand_id.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing AuctionResult".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(AuctionResult {
            auction_finish_ts: intermediate_rep.auction_finish_ts.into_iter().next(),
            auction_start_ts: intermediate_rep.auction_start_ts.into_iter().next(),
            banner: intermediate_rep.banner.into_iter().next(),
            bid_type: intermediate_rep.bid_type.into_iter().next(),
            interstitial: intermediate_rep.interstitial.into_iter().next(),
            price: intermediate_rep.price.into_iter().next(),
            rewarded: intermediate_rep.rewarded.into_iter().next(),
            status: intermediate_rep
                .status
                .into_iter()
                .next()
                .ok_or_else(|| "status missing in AuctionResult".to_string())?,
            winner_ad_unit_label: intermediate_rep.winner_ad_unit_label.into_iter().next(),
            winner_ad_unit_uid: intermediate_rep.winner_ad_unit_uid.into_iter().next(),
            winner_demand_id: intermediate_rep.winner_demand_id.into_iter().next(),
        })
    }
}

// Methods for converting between header::IntoHeaderValue<AuctionResult> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<AuctionResult>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<AuctionResult>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for AuctionResult - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<AuctionResult> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <AuctionResult as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into AuctionResult - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<AuctionResult>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<AuctionResult>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<AuctionResult>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<AuctionResult> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <AuctionResult as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into AuctionResult - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
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

impl std::fmt::Display for AuctionResultBidType {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match *self {
            AuctionResultBidType::Rtb => write!(f, "RTB"),
            AuctionResultBidType::Cpm => write!(f, "CPM"),
        }
    }
}

impl std::str::FromStr for AuctionResultBidType {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        match s {
            "RTB" => std::result::Result::Ok(AuctionResultBidType::Rtb),
            "CPM" => std::result::Result::Ok(AuctionResultBidType::Cpm),
            _ => std::result::Result::Err(format!("Value not valid: {}", s)),
        }
    }
}

// Methods for converting between header::IntoHeaderValue<AuctionResultBidType> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<AuctionResultBidType>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<AuctionResultBidType>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for AuctionResultBidType - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<AuctionResultBidType>
{
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <AuctionResultBidType as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into AuctionResultBidType - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<AuctionResultBidType>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<AuctionResultBidType>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<AuctionResultBidType>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values : std::vec::Vec<AuctionResultBidType> = hdr_values
                .split(',')
                .filter_map(|hdr_value| match hdr_value.trim() {
                    "" => std::option::Option::None,
                    hdr_value => std::option::Option::Some({
                        match <AuctionResultBidType as std::str::FromStr>::from_str(hdr_value) {
                            std::result::Result::Ok(value) => std::result::Result::Ok(value),
                            std::result::Result::Err(err) => std::result::Result::Err(
                                format!("Unable to convert header value '{}' into AuctionResultBidType - {}",
                                    hdr_value, err))
                        }
                    })
                }).collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
        }
    }
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

impl std::fmt::Display for AuctionResultStatus {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match *self {
            AuctionResultStatus::Success => write!(f, "SUCCESS"),
            AuctionResultStatus::Fail => write!(f, "FAIL"),
            AuctionResultStatus::AuctionCancelled => write!(f, "AUCTION_CANCELLED"),
        }
    }
}

impl std::str::FromStr for AuctionResultStatus {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        match s {
            "SUCCESS" => std::result::Result::Ok(AuctionResultStatus::Success),
            "FAIL" => std::result::Result::Ok(AuctionResultStatus::Fail),
            "AUCTION_CANCELLED" => std::result::Result::Ok(AuctionResultStatus::AuctionCancelled),
            _ => std::result::Result::Err(format!("Value not valid: {}", s)),
        }
    }
}

// Methods for converting between header::IntoHeaderValue<AuctionResultStatus> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<AuctionResultStatus>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<AuctionResultStatus>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for AuctionResultStatus - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<AuctionResultStatus>
{
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <AuctionResultStatus as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into AuctionResultStatus - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<AuctionResultStatus>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<AuctionResultStatus>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<AuctionResultStatus>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values : std::vec::Vec<AuctionResultStatus> = hdr_values
                .split(',')
                .filter_map(|hdr_value| match hdr_value.trim() {
                    "" => std::option::Option::None,
                    hdr_value => std::option::Option::Some({
                        match <AuctionResultStatus as std::str::FromStr>::from_str(hdr_value) {
                            std::result::Result::Ok(value) => std::result::Result::Ok(value),
                            std::result::Result::Err(err) => std::result::Result::Err(
                                format!("Unable to convert header value '{}' into AuctionResultStatus - {}",
                                    hdr_value, err))
                        }
                    })
                }).collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct BannerAdObject {
    #[serde(rename = "format")]
    pub format: models::AdFormat,
}

impl BannerAdObject {
    #[allow(clippy::new_without_default)]
    pub fn new(format: models::AdFormat) -> BannerAdObject {
        BannerAdObject { format }
    }
}

/// Converts the BannerAdObject value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for BannerAdObject {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            // Skipping non-primitive type format in query parameter serialization
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a BannerAdObject value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for BannerAdObject {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub format: Vec<models::AdFormat>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err(
                        "Missing value while parsing BannerAdObject".to_string(),
                    )
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    #[allow(clippy::redundant_clone)]
                    "format" => intermediate_rep.format.push(
                        <models::AdFormat as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing BannerAdObject".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(BannerAdObject {
            format: intermediate_rep
                .format
                .into_iter()
                .next()
                .ok_or_else(|| "format missing in BannerAdObject".to_string())?,
        })
    }
}

// Methods for converting between header::IntoHeaderValue<BannerAdObject> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<BannerAdObject>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<BannerAdObject>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for BannerAdObject - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<BannerAdObject> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <BannerAdObject as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into BannerAdObject - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<BannerAdObject>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<BannerAdObject>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<BannerAdObject>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<BannerAdObject> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <BannerAdObject as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into BannerAdObject - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct BaseRequest {
    #[serde(rename = "app")]
    pub app: models::App,

    #[serde(rename = "device")]
    pub device: models::Device,

    #[serde(rename = "ext")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ext: Option<String>,

    #[serde(rename = "geo")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub geo: Option<models::Geo>,

    #[serde(rename = "regs")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub regs: Option<models::Regulations>,

    #[serde(rename = "segment")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub segment: Option<models::Segment>,

    #[serde(rename = "session")]
    pub session: models::Session,

    #[serde(rename = "token")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub token: Option<String>,

    #[serde(rename = "user")]
    pub user: models::User,
}

impl BaseRequest {
    #[allow(clippy::new_without_default)]
    pub fn new(
        app: models::App,
        device: models::Device,
        session: models::Session,
        user: models::User,
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

/// Converts the BaseRequest value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for BaseRequest {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            // Skipping non-primitive type app in query parameter serialization
            // Skipping non-primitive type device in query parameter serialization
            self.ext
                .as_ref()
                .map(|ext| ["ext".to_string(), ext.to_string()].join(",")),
            // Skipping non-primitive type geo in query parameter serialization
            // Skipping non-primitive type regs in query parameter serialization
            // Skipping non-primitive type segment in query parameter serialization
            // Skipping non-primitive type session in query parameter serialization
            self.token
                .as_ref()
                .map(|token| ["token".to_string(), token.to_string()].join(",")),
            // Skipping non-primitive type user in query parameter serialization
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a BaseRequest value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for BaseRequest {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub app: Vec<models::App>,
            pub device: Vec<models::Device>,
            pub ext: Vec<String>,
            pub geo: Vec<models::Geo>,
            pub regs: Vec<models::Regulations>,
            pub segment: Vec<models::Segment>,
            pub session: Vec<models::Session>,
            pub token: Vec<String>,
            pub user: Vec<models::User>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err(
                        "Missing value while parsing BaseRequest".to_string(),
                    )
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    #[allow(clippy::redundant_clone)]
                    "app" => intermediate_rep.app.push(
                        <models::App as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "device" => intermediate_rep.device.push(
                        <models::Device as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "ext" => intermediate_rep.ext.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "geo" => intermediate_rep.geo.push(
                        <models::Geo as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "regs" => intermediate_rep.regs.push(
                        <models::Regulations as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "segment" => intermediate_rep.segment.push(
                        <models::Segment as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "session" => intermediate_rep.session.push(
                        <models::Session as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "token" => intermediate_rep.token.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "user" => intermediate_rep.user.push(
                        <models::User as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing BaseRequest".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(BaseRequest {
            app: intermediate_rep
                .app
                .into_iter()
                .next()
                .ok_or_else(|| "app missing in BaseRequest".to_string())?,
            device: intermediate_rep
                .device
                .into_iter()
                .next()
                .ok_or_else(|| "device missing in BaseRequest".to_string())?,
            ext: intermediate_rep.ext.into_iter().next(),
            geo: intermediate_rep.geo.into_iter().next(),
            regs: intermediate_rep.regs.into_iter().next(),
            segment: intermediate_rep.segment.into_iter().next(),
            session: intermediate_rep
                .session
                .into_iter()
                .next()
                .ok_or_else(|| "session missing in BaseRequest".to_string())?,
            token: intermediate_rep.token.into_iter().next(),
            user: intermediate_rep
                .user
                .into_iter()
                .next()
                .ok_or_else(|| "user missing in BaseRequest".to_string())?,
        })
    }
}

// Methods for converting between header::IntoHeaderValue<BaseRequest> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<BaseRequest>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<BaseRequest>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for BaseRequest - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<BaseRequest> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <BaseRequest as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into BaseRequest - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<BaseRequest>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<BaseRequest>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<BaseRequest>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<BaseRequest> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <BaseRequest as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into BaseRequest - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize)]
pub struct Bid();

/// Converts the Bid value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl ::std::string::ToString for Bid {
    fn to_string(&self) -> String {
        // ToString for this model is not supported
        "".to_string()
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a Bid value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl ::std::str::FromStr for Bid {
    type Err = &'static str;

    fn from_str(_s: &str) -> std::result::Result<Self, Self::Err> {
        std::result::Result::Err("Parsing Bid is not supported")
    }
}

// Methods for converting between header::IntoHeaderValue<Bid> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<Bid>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(hdr_value: header::IntoHeaderValue<Bid>) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for Bid - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<Bid> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => match <Bid as std::str::FromStr>::from_str(value) {
                std::result::Result::Ok(value) => {
                    std::result::Result::Ok(header::IntoHeaderValue(value))
                }
                std::result::Result::Err(err) => std::result::Result::Err(format!(
                    "Unable to convert header value '{}' into Bid - {}",
                    value, err
                )),
            },
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<Bid>>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<Bid>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<Vec<Bid>> {
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<Bid> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <Bid as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into Bid - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
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
pub enum BidType {
    #[serde(rename = "RTB")]
    Rtb,
    #[serde(rename = "CPM")]
    Cpm,
}

impl std::fmt::Display for BidType {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match *self {
            BidType::Rtb => write!(f, "RTB"),
            BidType::Cpm => write!(f, "CPM"),
        }
    }
}

impl std::str::FromStr for BidType {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        match s {
            "RTB" => std::result::Result::Ok(BidType::Rtb),
            "CPM" => std::result::Result::Ok(BidType::Cpm),
            _ => std::result::Result::Err(format!("Value not valid: {}", s)),
        }
    }
}

// Methods for converting between header::IntoHeaderValue<BidType> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<BidType>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<BidType>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for BidType - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<BidType> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <BidType as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into BidType - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<BidType>>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<BidType>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<Vec<BidType>> {
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<BidType> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <BidType as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into BidType - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
        }
    }
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

/// Converts the BiddingConfig value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for BiddingConfig {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> =
            vec![self.token_timeout_ms.as_ref().map(|token_timeout_ms| {
                ["token_timeout_ms".to_string(), token_timeout_ms.to_string()].join(",")
            })];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a BiddingConfig value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for BiddingConfig {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub token_timeout_ms: Vec<i64>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err(
                        "Missing value while parsing BiddingConfig".to_string(),
                    )
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    #[allow(clippy::redundant_clone)]
                    "token_timeout_ms" => intermediate_rep.token_timeout_ms.push(
                        <i64 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing BiddingConfig".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(BiddingConfig {
            token_timeout_ms: intermediate_rep.token_timeout_ms.into_iter().next(),
        })
    }
}

// Methods for converting between header::IntoHeaderValue<BiddingConfig> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<BiddingConfig>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<BiddingConfig>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for BiddingConfig - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<BiddingConfig> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <BiddingConfig as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into BiddingConfig - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<BiddingConfig>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<BiddingConfig>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<BiddingConfig>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<BiddingConfig> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <BiddingConfig as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into BiddingConfig - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct ClickRequest {
    #[serde(rename = "app")]
    pub app: models::App,

    #[serde(rename = "device")]
    pub device: models::Device,

    #[serde(rename = "ext")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ext: Option<String>,

    #[serde(rename = "geo")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub geo: Option<models::Geo>,

    #[serde(rename = "regs")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub regs: Option<models::Regulations>,

    #[serde(rename = "segment")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub segment: Option<models::Segment>,

    #[serde(rename = "session")]
    pub session: models::Session,

    #[serde(rename = "token")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub token: Option<String>,

    #[serde(rename = "user")]
    pub user: models::User,

    #[serde(rename = "bid")]
    #[serde(deserialize_with = "swagger::nullable_format::deserialize_optional_nullable")]
    #[serde(default = "swagger::nullable_format::default_optional_nullable")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub bid: Option<swagger::Nullable<models::Bid>>,

    #[serde(rename = "show")]
    #[serde(deserialize_with = "swagger::nullable_format::deserialize_optional_nullable")]
    #[serde(default = "swagger::nullable_format::default_optional_nullable")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub show: Option<swagger::Nullable<models::Bid>>,
}

impl ClickRequest {
    #[allow(clippy::new_without_default)]
    pub fn new(
        app: models::App,
        device: models::Device,
        session: models::Session,
        user: models::User,
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

/// Converts the ClickRequest value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for ClickRequest {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            // Skipping non-primitive type app in query parameter serialization
            // Skipping non-primitive type device in query parameter serialization
            self.ext
                .as_ref()
                .map(|ext| ["ext".to_string(), ext.to_string()].join(",")),
            // Skipping non-primitive type geo in query parameter serialization
            // Skipping non-primitive type regs in query parameter serialization
            // Skipping non-primitive type segment in query parameter serialization
            // Skipping non-primitive type session in query parameter serialization
            self.token
                .as_ref()
                .map(|token| ["token".to_string(), token.to_string()].join(",")),
            // Skipping non-primitive type user in query parameter serialization
            // Skipping non-primitive type bid in query parameter serialization
            // Skipping non-primitive type show in query parameter serialization
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a ClickRequest value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for ClickRequest {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub app: Vec<models::App>,
            pub device: Vec<models::Device>,
            pub ext: Vec<String>,
            pub geo: Vec<models::Geo>,
            pub regs: Vec<models::Regulations>,
            pub segment: Vec<models::Segment>,
            pub session: Vec<models::Session>,
            pub token: Vec<String>,
            pub user: Vec<models::User>,
            pub bid: Vec<swagger::Nullable<models::Bid>>,
            pub show: Vec<swagger::Nullable<models::Bid>>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err(
                        "Missing value while parsing ClickRequest".to_string(),
                    )
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    #[allow(clippy::redundant_clone)]
                    "app" => intermediate_rep.app.push(
                        <models::App as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "device" => intermediate_rep.device.push(
                        <models::Device as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "ext" => intermediate_rep.ext.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "geo" => intermediate_rep.geo.push(
                        <models::Geo as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "regs" => intermediate_rep.regs.push(
                        <models::Regulations as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "segment" => intermediate_rep.segment.push(
                        <models::Segment as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "session" => intermediate_rep.session.push(
                        <models::Session as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "token" => intermediate_rep.token.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "user" => intermediate_rep.user.push(
                        <models::User as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    "bid" => return std::result::Result::Err(
                        "Parsing a nullable type in this style is not supported in ClickRequest"
                            .to_string(),
                    ),
                    "show" => return std::result::Result::Err(
                        "Parsing a nullable type in this style is not supported in ClickRequest"
                            .to_string(),
                    ),
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing ClickRequest".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(ClickRequest {
            app: intermediate_rep
                .app
                .into_iter()
                .next()
                .ok_or_else(|| "app missing in ClickRequest".to_string())?,
            device: intermediate_rep
                .device
                .into_iter()
                .next()
                .ok_or_else(|| "device missing in ClickRequest".to_string())?,
            ext: intermediate_rep.ext.into_iter().next(),
            geo: intermediate_rep.geo.into_iter().next(),
            regs: intermediate_rep.regs.into_iter().next(),
            segment: intermediate_rep.segment.into_iter().next(),
            session: intermediate_rep
                .session
                .into_iter()
                .next()
                .ok_or_else(|| "session missing in ClickRequest".to_string())?,
            token: intermediate_rep.token.into_iter().next(),
            user: intermediate_rep
                .user
                .into_iter()
                .next()
                .ok_or_else(|| "user missing in ClickRequest".to_string())?,
            bid: std::result::Result::Err(
                "Nullable types not supported in ClickRequest".to_string(),
            )?,
            show: std::result::Result::Err(
                "Nullable types not supported in ClickRequest".to_string(),
            )?,
        })
    }
}

// Methods for converting between header::IntoHeaderValue<ClickRequest> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<ClickRequest>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<ClickRequest>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for ClickRequest - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<ClickRequest> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <ClickRequest as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into ClickRequest - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<ClickRequest>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<ClickRequest>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<ClickRequest>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<ClickRequest> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <ClickRequest as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into ClickRequest - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct ConfigRequest {
    #[serde(rename = "app")]
    pub app: models::App,

    #[serde(rename = "device")]
    pub device: models::Device,

    #[serde(rename = "ext")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ext: Option<String>,

    #[serde(rename = "geo")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub geo: Option<models::Geo>,

    #[serde(rename = "regs")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub regs: Option<models::Regulations>,

    #[serde(rename = "segment")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub segment: Option<models::Segment>,

    #[serde(rename = "session")]
    pub session: models::Session,

    #[serde(rename = "token")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub token: Option<String>,

    #[serde(rename = "user")]
    pub user: models::User,

    #[serde(rename = "adapters")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub adapters: Option<std::collections::HashMap<String, models::Adapter>>,
}

impl ConfigRequest {
    #[allow(clippy::new_without_default)]
    pub fn new(
        app: models::App,
        device: models::Device,
        session: models::Session,
        user: models::User,
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

/// Converts the ConfigRequest value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for ConfigRequest {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            // Skipping non-primitive type app in query parameter serialization
            // Skipping non-primitive type device in query parameter serialization
            self.ext
                .as_ref()
                .map(|ext| ["ext".to_string(), ext.to_string()].join(",")),
            // Skipping non-primitive type geo in query parameter serialization
            // Skipping non-primitive type regs in query parameter serialization
            // Skipping non-primitive type segment in query parameter serialization
            // Skipping non-primitive type session in query parameter serialization
            self.token
                .as_ref()
                .map(|token| ["token".to_string(), token.to_string()].join(",")),
            // Skipping non-primitive type user in query parameter serialization
            // Skipping map adapters in query parameter serialization
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a ConfigRequest value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for ConfigRequest {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub app: Vec<models::App>,
            pub device: Vec<models::Device>,
            pub ext: Vec<String>,
            pub geo: Vec<models::Geo>,
            pub regs: Vec<models::Regulations>,
            pub segment: Vec<models::Segment>,
            pub session: Vec<models::Session>,
            pub token: Vec<String>,
            pub user: Vec<models::User>,
            pub adapters: Vec<std::collections::HashMap<String, models::Adapter>>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err(
                        "Missing value while parsing ConfigRequest".to_string(),
                    )
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    #[allow(clippy::redundant_clone)]
                    "app" => intermediate_rep.app.push(
                        <models::App as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "device" => intermediate_rep.device.push(
                        <models::Device as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "ext" => intermediate_rep.ext.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "geo" => intermediate_rep.geo.push(
                        <models::Geo as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "regs" => intermediate_rep.regs.push(
                        <models::Regulations as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "segment" => intermediate_rep.segment.push(
                        <models::Segment as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "session" => intermediate_rep.session.push(
                        <models::Session as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "token" => intermediate_rep.token.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "user" => intermediate_rep.user.push(
                        <models::User as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    "adapters" => {
                        return std::result::Result::Err(
                            "Parsing a container in this style is not supported in ConfigRequest"
                                .to_string(),
                        )
                    }
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing ConfigRequest".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(ConfigRequest {
            app: intermediate_rep
                .app
                .into_iter()
                .next()
                .ok_or_else(|| "app missing in ConfigRequest".to_string())?,
            device: intermediate_rep
                .device
                .into_iter()
                .next()
                .ok_or_else(|| "device missing in ConfigRequest".to_string())?,
            ext: intermediate_rep.ext.into_iter().next(),
            geo: intermediate_rep.geo.into_iter().next(),
            regs: intermediate_rep.regs.into_iter().next(),
            segment: intermediate_rep.segment.into_iter().next(),
            session: intermediate_rep
                .session
                .into_iter()
                .next()
                .ok_or_else(|| "session missing in ConfigRequest".to_string())?,
            token: intermediate_rep.token.into_iter().next(),
            user: intermediate_rep
                .user
                .into_iter()
                .next()
                .ok_or_else(|| "user missing in ConfigRequest".to_string())?,
            adapters: intermediate_rep.adapters.into_iter().next(),
        })
    }
}

// Methods for converting between header::IntoHeaderValue<ConfigRequest> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<ConfigRequest>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<ConfigRequest>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for ConfigRequest - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<ConfigRequest> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <ConfigRequest as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into ConfigRequest - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<ConfigRequest>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<ConfigRequest>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<ConfigRequest>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<ConfigRequest> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <ConfigRequest as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into ConfigRequest - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct ConfigResponse {
    #[serde(rename = "bidding")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub bidding: Option<models::BiddingConfig>,

    #[serde(rename = "init")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub init: Option<models::InitConfig>,

    #[serde(rename = "placements")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub placements: Option<Vec<serde_json::Value>>,

    #[serde(rename = "segment")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub segment: Option<models::Segment>,

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

/// Converts the ConfigResponse value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for ConfigResponse {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            // Skipping non-primitive type bidding in query parameter serialization
            // Skipping non-primitive type init in query parameter serialization
            // Skipping non-primitive type placements in query parameter serialization
            // Skipping non-primitive type segment in query parameter serialization
            self.token
                .as_ref()
                .map(|token| ["token".to_string(), token.to_string()].join(",")),
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a ConfigResponse value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for ConfigResponse {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub bidding: Vec<models::BiddingConfig>,
            pub init: Vec<models::InitConfig>,
            pub placements: Vec<Vec<serde_json::Value>>,
            pub segment: Vec<models::Segment>,
            pub token: Vec<String>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err(
                        "Missing value while parsing ConfigResponse".to_string(),
                    )
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    #[allow(clippy::redundant_clone)]
                    "bidding" => intermediate_rep.bidding.push(
                        <models::BiddingConfig as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "init" => intermediate_rep.init.push(
                        <models::InitConfig as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    "placements" => {
                        return std::result::Result::Err(
                            "Parsing a container in this style is not supported in ConfigResponse"
                                .to_string(),
                        )
                    }
                    #[allow(clippy::redundant_clone)]
                    "segment" => intermediate_rep.segment.push(
                        <models::Segment as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "token" => intermediate_rep.token.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing ConfigResponse".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(ConfigResponse {
            bidding: intermediate_rep.bidding.into_iter().next(),
            init: intermediate_rep.init.into_iter().next(),
            placements: intermediate_rep.placements.into_iter().next(),
            segment: intermediate_rep.segment.into_iter().next(),
            token: intermediate_rep.token.into_iter().next(),
        })
    }
}

// Methods for converting between header::IntoHeaderValue<ConfigResponse> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<ConfigResponse>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<ConfigResponse>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for ConfigResponse - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<ConfigResponse> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <ConfigResponse as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into ConfigResponse - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<ConfigResponse>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<ConfigResponse>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<ConfigResponse>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<ConfigResponse> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <ConfigResponse as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into ConfigResponse - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
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
    pub connection_type: models::DeviceConnectionType,

    #[serde(rename = "geo")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub geo: Option<models::Geo>,

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
    pub r#type: Option<models::DeviceType>,

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
        connection_type: models::DeviceConnectionType,
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
            r#type: None,
            ua,
            w,
        }
    }
}

/// Converts the Device value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for Device {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            self.carrier
                .as_ref()
                .map(|carrier| ["carrier".to_string(), carrier.to_string()].join(",")),
            // Skipping non-primitive type connection_type in query parameter serialization
            // Skipping non-primitive type geo in query parameter serialization
            Some("h".to_string()),
            Some(self.h.to_string()),
            Some("hwv".to_string()),
            Some(self.hwv.to_string()),
            Some("js".to_string()),
            Some(self.js.to_string()),
            Some("language".to_string()),
            Some(self.language.to_string()),
            Some("make".to_string()),
            Some(self.make.to_string()),
            self.mccmnc
                .as_ref()
                .map(|mccmnc| ["mccmnc".to_string(), mccmnc.to_string()].join(",")),
            Some("model".to_string()),
            Some(self.model.to_string()),
            Some("os".to_string()),
            Some(self.os.to_string()),
            Some("osv".to_string()),
            Some(self.osv.to_string()),
            Some("ppi".to_string()),
            Some(self.ppi.to_string()),
            Some("pxratio".to_string()),
            Some(self.pxratio.to_string()),
            // Skipping non-primitive type type in query parameter serialization
            Some("ua".to_string()),
            Some(self.ua.to_string()),
            Some("w".to_string()),
            Some(self.w.to_string()),
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a Device value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for Device {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub carrier: Vec<String>,
            pub connection_type: Vec<models::DeviceConnectionType>,
            pub geo: Vec<models::Geo>,
            pub h: Vec<i32>,
            pub hwv: Vec<String>,
            pub js: Vec<i32>,
            pub language: Vec<String>,
            pub make: Vec<String>,
            pub mccmnc: Vec<String>,
            pub model: Vec<String>,
            pub os: Vec<String>,
            pub osv: Vec<String>,
            pub ppi: Vec<i32>,
            pub pxratio: Vec<f64>,
            pub r#type: Vec<models::DeviceType>,
            pub ua: Vec<String>,
            pub w: Vec<i32>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err(
                        "Missing value while parsing Device".to_string(),
                    )
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    #[allow(clippy::redundant_clone)]
                    "carrier" => intermediate_rep.carrier.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "connection_type" => intermediate_rep.connection_type.push(
                        <models::DeviceConnectionType as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "geo" => intermediate_rep.geo.push(
                        <models::Geo as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "h" => intermediate_rep.h.push(
                        <i32 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "hwv" => intermediate_rep.hwv.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "js" => intermediate_rep.js.push(
                        <i32 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "language" => intermediate_rep.language.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "make" => intermediate_rep.make.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "mccmnc" => intermediate_rep.mccmnc.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "model" => intermediate_rep.model.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "os" => intermediate_rep.os.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "osv" => intermediate_rep.osv.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "ppi" => intermediate_rep.ppi.push(
                        <i32 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "pxratio" => intermediate_rep.pxratio.push(
                        <f64 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "type" => intermediate_rep.r#type.push(
                        <models::DeviceType as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "ua" => intermediate_rep.ua.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "w" => intermediate_rep.w.push(
                        <i32 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing Device".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(Device {
            carrier: intermediate_rep.carrier.into_iter().next(),
            connection_type: intermediate_rep
                .connection_type
                .into_iter()
                .next()
                .ok_or_else(|| "connection_type missing in Device".to_string())?,
            geo: intermediate_rep.geo.into_iter().next(),
            h: intermediate_rep
                .h
                .into_iter()
                .next()
                .ok_or_else(|| "h missing in Device".to_string())?,
            hwv: intermediate_rep
                .hwv
                .into_iter()
                .next()
                .ok_or_else(|| "hwv missing in Device".to_string())?,
            js: intermediate_rep
                .js
                .into_iter()
                .next()
                .ok_or_else(|| "js missing in Device".to_string())?,
            language: intermediate_rep
                .language
                .into_iter()
                .next()
                .ok_or_else(|| "language missing in Device".to_string())?,
            make: intermediate_rep
                .make
                .into_iter()
                .next()
                .ok_or_else(|| "make missing in Device".to_string())?,
            mccmnc: intermediate_rep.mccmnc.into_iter().next(),
            model: intermediate_rep
                .model
                .into_iter()
                .next()
                .ok_or_else(|| "model missing in Device".to_string())?,
            os: intermediate_rep
                .os
                .into_iter()
                .next()
                .ok_or_else(|| "os missing in Device".to_string())?,
            osv: intermediate_rep
                .osv
                .into_iter()
                .next()
                .ok_or_else(|| "osv missing in Device".to_string())?,
            ppi: intermediate_rep
                .ppi
                .into_iter()
                .next()
                .ok_or_else(|| "ppi missing in Device".to_string())?,
            pxratio: intermediate_rep
                .pxratio
                .into_iter()
                .next()
                .ok_or_else(|| "pxratio missing in Device".to_string())?,
            r#type: intermediate_rep.r#type.into_iter().next(),
            ua: intermediate_rep
                .ua
                .into_iter()
                .next()
                .ok_or_else(|| "ua missing in Device".to_string())?,
            w: intermediate_rep
                .w
                .into_iter()
                .next()
                .ok_or_else(|| "w missing in Device".to_string())?,
        })
    }
}

// Methods for converting between header::IntoHeaderValue<Device> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<Device>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<Device>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for Device - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<Device> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <Device as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into Device - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<Device>>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<Device>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<Vec<Device>> {
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<Device> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <Device as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into Device - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
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

impl std::fmt::Display for DeviceConnectionType {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match *self {
            DeviceConnectionType::Ethernet => write!(f, "ETHERNET"),
            DeviceConnectionType::Wifi => write!(f, "WIFI"),
            DeviceConnectionType::Cellular => write!(f, "CELLULAR"),
            DeviceConnectionType::CellularUnknown => write!(f, "CELLULAR_UNKNOWN"),
            DeviceConnectionType::Cellular2G => write!(f, "CELLULAR_2_G"),
            DeviceConnectionType::Cellular3G => write!(f, "CELLULAR_3_G"),
            DeviceConnectionType::Cellular4G => write!(f, "CELLULAR_4_G"),
            DeviceConnectionType::Cellular5G => write!(f, "CELLULAR_5_G"),
        }
    }
}

impl std::str::FromStr for DeviceConnectionType {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        match s {
            "ETHERNET" => std::result::Result::Ok(DeviceConnectionType::Ethernet),
            "WIFI" => std::result::Result::Ok(DeviceConnectionType::Wifi),
            "CELLULAR" => std::result::Result::Ok(DeviceConnectionType::Cellular),
            "CELLULAR_UNKNOWN" => std::result::Result::Ok(DeviceConnectionType::CellularUnknown),
            "CELLULAR_2_G" => std::result::Result::Ok(DeviceConnectionType::Cellular2G),
            "CELLULAR_3_G" => std::result::Result::Ok(DeviceConnectionType::Cellular3G),
            "CELLULAR_4_G" => std::result::Result::Ok(DeviceConnectionType::Cellular4G),
            "CELLULAR_5_G" => std::result::Result::Ok(DeviceConnectionType::Cellular5G),
            _ => std::result::Result::Err(format!("Value not valid: {}", s)),
        }
    }
}

// Methods for converting between header::IntoHeaderValue<DeviceConnectionType> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<DeviceConnectionType>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<DeviceConnectionType>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for DeviceConnectionType - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<DeviceConnectionType>
{
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <DeviceConnectionType as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into DeviceConnectionType - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<DeviceConnectionType>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<DeviceConnectionType>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<DeviceConnectionType>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values : std::vec::Vec<DeviceConnectionType> = hdr_values
                .split(',')
                .filter_map(|hdr_value| match hdr_value.trim() {
                    "" => std::option::Option::None,
                    hdr_value => std::option::Option::Some({
                        match <DeviceConnectionType as std::str::FromStr>::from_str(hdr_value) {
                            std::result::Result::Ok(value) => std::result::Result::Ok(value),
                            std::result::Result::Err(err) => std::result::Result::Err(
                                format!("Unable to convert header value '{}' into DeviceConnectionType - {}",
                                    hdr_value, err))
                        }
                    })
                }).collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
        }
    }
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

impl std::fmt::Display for DeviceType {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match *self {
            DeviceType::Phone => write!(f, "PHONE"),
            DeviceType::Tablet => write!(f, "TABLET"),
        }
    }
}

impl std::str::FromStr for DeviceType {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        match s {
            "PHONE" => std::result::Result::Ok(DeviceType::Phone),
            "TABLET" => std::result::Result::Ok(DeviceType::Tablet),
            _ => std::result::Result::Err(format!("Value not valid: {}", s)),
        }
    }
}

// Methods for converting between header::IntoHeaderValue<DeviceType> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<DeviceType>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<DeviceType>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for DeviceType - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<DeviceType> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <DeviceType as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into DeviceType - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<DeviceType>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<DeviceType>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<DeviceType>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<DeviceType> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <DeviceType as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into DeviceType - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct Error {
    #[serde(rename = "error")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub error: Option<models::ErrorError>,
}

impl Error {
    #[allow(clippy::new_without_default)]
    pub fn new() -> Error {
        Error { error: None }
    }
}

/// Converts the Error value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for Error {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            // Skipping non-primitive type error in query parameter serialization
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a Error value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for Error {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub error: Vec<models::ErrorError>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err(
                        "Missing value while parsing Error".to_string(),
                    )
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    #[allow(clippy::redundant_clone)]
                    "error" => intermediate_rep.error.push(
                        <models::ErrorError as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing Error".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(Error {
            error: intermediate_rep.error.into_iter().next(),
        })
    }
}

// Methods for converting between header::IntoHeaderValue<Error> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<Error>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<Error>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for Error - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<Error> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => match <Error as std::str::FromStr>::from_str(value) {
                std::result::Result::Ok(value) => {
                    std::result::Result::Ok(header::IntoHeaderValue(value))
                }
                std::result::Result::Err(err) => std::result::Result::Err(format!(
                    "Unable to convert header value '{}' into Error - {}",
                    value, err
                )),
            },
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<Error>>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<Error>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<Vec<Error>> {
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<Error> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <Error as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into Error - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
        }
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

/// Converts the ErrorError value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for ErrorError {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            self.code
                .as_ref()
                .map(|code| ["code".to_string(), code.to_string()].join(",")),
            self.message
                .as_ref()
                .map(|message| ["message".to_string(), message.to_string()].join(",")),
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a ErrorError value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for ErrorError {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub code: Vec<i32>,
            pub message: Vec<String>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err(
                        "Missing value while parsing ErrorError".to_string(),
                    )
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    #[allow(clippy::redundant_clone)]
                    "code" => intermediate_rep.code.push(
                        <i32 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "message" => intermediate_rep.message.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing ErrorError".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(ErrorError {
            code: intermediate_rep.code.into_iter().next(),
            message: intermediate_rep.message.into_iter().next(),
        })
    }
}

// Methods for converting between header::IntoHeaderValue<ErrorError> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<ErrorError>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<ErrorError>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for ErrorError - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<ErrorError> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <ErrorError as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into ErrorError - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<ErrorError>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<ErrorError>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<ErrorError>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<ErrorError> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <ErrorError as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into ErrorError - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
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

/// Converts the ExternalWinner value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for ExternalWinner {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            self.demand_id
                .as_ref()
                .map(|demand_id| ["demand_id".to_string(), demand_id.to_string()].join(",")),
            Some("price".to_string()),
            Some(self.price.to_string()),
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a ExternalWinner value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for ExternalWinner {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub demand_id: Vec<String>,
            pub price: Vec<f64>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err(
                        "Missing value while parsing ExternalWinner".to_string(),
                    )
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    #[allow(clippy::redundant_clone)]
                    "demand_id" => intermediate_rep.demand_id.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "price" => intermediate_rep.price.push(
                        <f64 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing ExternalWinner".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(ExternalWinner {
            demand_id: intermediate_rep.demand_id.into_iter().next(),
            price: intermediate_rep
                .price
                .into_iter()
                .next()
                .ok_or_else(|| "price missing in ExternalWinner".to_string())?,
        })
    }
}

// Methods for converting between header::IntoHeaderValue<ExternalWinner> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<ExternalWinner>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<ExternalWinner>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for ExternalWinner - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<ExternalWinner> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <ExternalWinner as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into ExternalWinner - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<ExternalWinner>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<ExternalWinner>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<ExternalWinner>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<ExternalWinner> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <ExternalWinner as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into ExternalWinner - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
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

/// Converts the Geo value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for Geo {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            self.accuracy
                .as_ref()
                .map(|accuracy| ["accuracy".to_string(), accuracy.to_string()].join(",")),
            self.city
                .as_ref()
                .map(|city| ["city".to_string(), city.to_string()].join(",")),
            self.country
                .as_ref()
                .map(|country| ["country".to_string(), country.to_string()].join(",")),
            self.lastfix
                .as_ref()
                .map(|lastfix| ["lastfix".to_string(), lastfix.to_string()].join(",")),
            self.lat
                .as_ref()
                .map(|lat| ["lat".to_string(), lat.to_string()].join(",")),
            self.lon
                .as_ref()
                .map(|lon| ["lon".to_string(), lon.to_string()].join(",")),
            self.utcoffset
                .as_ref()
                .map(|utcoffset| ["utcoffset".to_string(), utcoffset.to_string()].join(",")),
            self.zip
                .as_ref()
                .map(|zip| ["zip".to_string(), zip.to_string()].join(",")),
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a Geo value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for Geo {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub accuracy: Vec<f64>,
            pub city: Vec<String>,
            pub country: Vec<String>,
            pub lastfix: Vec<i32>,
            pub lat: Vec<f64>,
            pub lon: Vec<f64>,
            pub utcoffset: Vec<i32>,
            pub zip: Vec<String>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err("Missing value while parsing Geo".to_string())
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    #[allow(clippy::redundant_clone)]
                    "accuracy" => intermediate_rep.accuracy.push(
                        <f64 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "city" => intermediate_rep.city.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "country" => intermediate_rep.country.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "lastfix" => intermediate_rep.lastfix.push(
                        <i32 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "lat" => intermediate_rep.lat.push(
                        <f64 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "lon" => intermediate_rep.lon.push(
                        <f64 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "utcoffset" => intermediate_rep.utcoffset.push(
                        <i32 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "zip" => intermediate_rep.zip.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing Geo".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(Geo {
            accuracy: intermediate_rep.accuracy.into_iter().next(),
            city: intermediate_rep.city.into_iter().next(),
            country: intermediate_rep.country.into_iter().next(),
            lastfix: intermediate_rep.lastfix.into_iter().next(),
            lat: intermediate_rep.lat.into_iter().next(),
            lon: intermediate_rep.lon.into_iter().next(),
            utcoffset: intermediate_rep.utcoffset.into_iter().next(),
            zip: intermediate_rep.zip.into_iter().next(),
        })
    }
}

// Methods for converting between header::IntoHeaderValue<Geo> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<Geo>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(hdr_value: header::IntoHeaderValue<Geo>) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for Geo - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<Geo> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => match <Geo as std::str::FromStr>::from_str(value) {
                std::result::Result::Ok(value) => {
                    std::result::Result::Ok(header::IntoHeaderValue(value))
                }
                std::result::Result::Err(err) => std::result::Result::Err(format!(
                    "Unable to convert header value '{}' into Geo - {}",
                    value, err
                )),
            },
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<Geo>>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<Geo>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<Vec<Geo>> {
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<Geo> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <Geo as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into Geo - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
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

// Methods for converting between header::IntoHeaderValue<GetAuctionAdTypeParameter> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<GetAuctionAdTypeParameter>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<GetAuctionAdTypeParameter>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for GetAuctionAdTypeParameter - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<GetAuctionAdTypeParameter>
{
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <GetAuctionAdTypeParameter as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into GetAuctionAdTypeParameter - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<GetAuctionAdTypeParameter>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<GetAuctionAdTypeParameter>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<GetAuctionAdTypeParameter>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values : std::vec::Vec<GetAuctionAdTypeParameter> = hdr_values
                .split(',')
                .filter_map(|hdr_value| match hdr_value.trim() {
                    "" => std::option::Option::None,
                    hdr_value => std::option::Option::Some({
                        match <GetAuctionAdTypeParameter as std::str::FromStr>::from_str(hdr_value) {
                            std::result::Result::Ok(value) => std::result::Result::Ok(value),
                            std::result::Result::Err(err) => std::result::Result::Err(
                                format!("Unable to convert header value '{}' into GetAuctionAdTypeParameter - {}",
                                    hdr_value, err))
                        }
                    })
                }).collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
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

/// Converts the InitConfig value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for InitConfig {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            // Skipping map adapters in query parameter serialization
            Some("tmax".to_string()),
            Some(self.tmax.to_string()),
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a InitConfig value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for InitConfig {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub adapters: Vec<std::collections::HashMap<String, serde_json::Value>>,
            pub tmax: Vec<i64>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err(
                        "Missing value while parsing InitConfig".to_string(),
                    )
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    "adapters" => {
                        return std::result::Result::Err(
                            "Parsing a container in this style is not supported in InitConfig"
                                .to_string(),
                        )
                    }
                    #[allow(clippy::redundant_clone)]
                    "tmax" => intermediate_rep.tmax.push(
                        <i64 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing InitConfig".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(InitConfig {
            adapters: intermediate_rep
                .adapters
                .into_iter()
                .next()
                .ok_or_else(|| "adapters missing in InitConfig".to_string())?,
            tmax: intermediate_rep
                .tmax
                .into_iter()
                .next()
                .ok_or_else(|| "tmax missing in InitConfig".to_string())?,
        })
    }
}

// Methods for converting between header::IntoHeaderValue<InitConfig> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<InitConfig>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<InitConfig>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for InitConfig - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<InitConfig> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <InitConfig as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into InitConfig - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<InitConfig>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<InitConfig>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<InitConfig>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<InitConfig> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <InitConfig as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into InitConfig - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct LossRequest {
    #[serde(rename = "app")]
    pub app: models::App,

    #[serde(rename = "device")]
    pub device: models::Device,

    #[serde(rename = "ext")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ext: Option<String>,

    #[serde(rename = "geo")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub geo: Option<models::Geo>,

    #[serde(rename = "regs")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub regs: Option<models::Regulations>,

    #[serde(rename = "segment")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub segment: Option<models::Segment>,

    #[serde(rename = "session")]
    pub session: models::Session,

    #[serde(rename = "token")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub token: Option<String>,

    #[serde(rename = "user")]
    pub user: models::User,

    #[serde(rename = "bid")]
    #[serde(deserialize_with = "swagger::nullable_format::deserialize_optional_nullable")]
    #[serde(default = "swagger::nullable_format::default_optional_nullable")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub bid: Option<swagger::Nullable<models::Bid>>,

    #[serde(rename = "show")]
    #[serde(deserialize_with = "swagger::nullable_format::deserialize_optional_nullable")]
    #[serde(default = "swagger::nullable_format::default_optional_nullable")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub show: Option<swagger::Nullable<models::Bid>>,

    #[serde(rename = "external_winner")]
    pub external_winner: models::ExternalWinner,
}

impl LossRequest {
    #[allow(clippy::new_without_default)]
    pub fn new(
        app: models::App,
        device: models::Device,
        session: models::Session,
        user: models::User,
        external_winner: models::ExternalWinner,
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

/// Converts the LossRequest value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for LossRequest {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            // Skipping non-primitive type app in query parameter serialization
            // Skipping non-primitive type device in query parameter serialization
            self.ext
                .as_ref()
                .map(|ext| ["ext".to_string(), ext.to_string()].join(",")),
            // Skipping non-primitive type geo in query parameter serialization
            // Skipping non-primitive type regs in query parameter serialization
            // Skipping non-primitive type segment in query parameter serialization
            // Skipping non-primitive type session in query parameter serialization
            self.token
                .as_ref()
                .map(|token| ["token".to_string(), token.to_string()].join(",")),
            // Skipping non-primitive type user in query parameter serialization
            // Skipping non-primitive type bid in query parameter serialization
            // Skipping non-primitive type show in query parameter serialization
            // Skipping non-primitive type external_winner in query parameter serialization
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a LossRequest value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for LossRequest {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub app: Vec<models::App>,
            pub device: Vec<models::Device>,
            pub ext: Vec<String>,
            pub geo: Vec<models::Geo>,
            pub regs: Vec<models::Regulations>,
            pub segment: Vec<models::Segment>,
            pub session: Vec<models::Session>,
            pub token: Vec<String>,
            pub user: Vec<models::User>,
            pub bid: Vec<swagger::Nullable<models::Bid>>,
            pub show: Vec<swagger::Nullable<models::Bid>>,
            pub external_winner: Vec<models::ExternalWinner>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err(
                        "Missing value while parsing LossRequest".to_string(),
                    )
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    #[allow(clippy::redundant_clone)]
                    "app" => intermediate_rep.app.push(
                        <models::App as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "device" => intermediate_rep.device.push(
                        <models::Device as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "ext" => intermediate_rep.ext.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "geo" => intermediate_rep.geo.push(
                        <models::Geo as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "regs" => intermediate_rep.regs.push(
                        <models::Regulations as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "segment" => intermediate_rep.segment.push(
                        <models::Segment as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "session" => intermediate_rep.session.push(
                        <models::Session as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "token" => intermediate_rep.token.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "user" => intermediate_rep.user.push(
                        <models::User as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    "bid" => {
                        return std::result::Result::Err(
                            "Parsing a nullable type in this style is not supported in LossRequest"
                                .to_string(),
                        )
                    }
                    "show" => {
                        return std::result::Result::Err(
                            "Parsing a nullable type in this style is not supported in LossRequest"
                                .to_string(),
                        )
                    }
                    #[allow(clippy::redundant_clone)]
                    "external_winner" => intermediate_rep.external_winner.push(
                        <models::ExternalWinner as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing LossRequest".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(LossRequest {
            app: intermediate_rep
                .app
                .into_iter()
                .next()
                .ok_or_else(|| "app missing in LossRequest".to_string())?,
            device: intermediate_rep
                .device
                .into_iter()
                .next()
                .ok_or_else(|| "device missing in LossRequest".to_string())?,
            ext: intermediate_rep.ext.into_iter().next(),
            geo: intermediate_rep.geo.into_iter().next(),
            regs: intermediate_rep.regs.into_iter().next(),
            segment: intermediate_rep.segment.into_iter().next(),
            session: intermediate_rep
                .session
                .into_iter()
                .next()
                .ok_or_else(|| "session missing in LossRequest".to_string())?,
            token: intermediate_rep.token.into_iter().next(),
            user: intermediate_rep
                .user
                .into_iter()
                .next()
                .ok_or_else(|| "user missing in LossRequest".to_string())?,
            bid: std::result::Result::Err(
                "Nullable types not supported in LossRequest".to_string(),
            )?,
            show: std::result::Result::Err(
                "Nullable types not supported in LossRequest".to_string(),
            )?,
            external_winner: intermediate_rep
                .external_winner
                .into_iter()
                .next()
                .ok_or_else(|| "external_winner missing in LossRequest".to_string())?,
        })
    }
}

// Methods for converting between header::IntoHeaderValue<LossRequest> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<LossRequest>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<LossRequest>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for LossRequest - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<LossRequest> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <LossRequest as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into LossRequest - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<LossRequest>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<LossRequest>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<LossRequest>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<LossRequest> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <LossRequest as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into LossRequest - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
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

impl std::fmt::Display for PostRewardAdTypeParameter {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match *self {
            PostRewardAdTypeParameter::Rewarded => write!(f, "rewarded"),
        }
    }
}

impl std::str::FromStr for PostRewardAdTypeParameter {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        match s {
            "rewarded" => std::result::Result::Ok(PostRewardAdTypeParameter::Rewarded),
            _ => std::result::Result::Err(format!("Value not valid: {}", s)),
        }
    }
}

// Methods for converting between header::IntoHeaderValue<PostRewardAdTypeParameter> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<PostRewardAdTypeParameter>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<PostRewardAdTypeParameter>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for PostRewardAdTypeParameter - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<PostRewardAdTypeParameter>
{
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <PostRewardAdTypeParameter as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into PostRewardAdTypeParameter - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<PostRewardAdTypeParameter>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<PostRewardAdTypeParameter>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<PostRewardAdTypeParameter>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values : std::vec::Vec<PostRewardAdTypeParameter> = hdr_values
                .split(',')
                .filter_map(|hdr_value| match hdr_value.trim() {
                    "" => std::option::Option::None,
                    hdr_value => std::option::Option::Some({
                        match <PostRewardAdTypeParameter as std::str::FromStr>::from_str(hdr_value) {
                            std::result::Result::Ok(value) => std::result::Result::Ok(value),
                            std::result::Result::Err(err) => std::result::Result::Err(
                                format!("Unable to convert header value '{}' into PostRewardAdTypeParameter - {}",
                                    hdr_value, err))
                        }
                    })
                }).collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
        }
    }
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

/// Converts the Regulations value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for Regulations {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            self.coppa
                .as_ref()
                .map(|coppa| ["coppa".to_string(), coppa.to_string()].join(",")),
            self.eu_privacy
                .as_ref()
                .map(|eu_privacy| ["eu_privacy".to_string(), eu_privacy.to_string()].join(",")),
            self.gdpr
                .as_ref()
                .map(|gdpr| ["gdpr".to_string(), gdpr.to_string()].join(",")),
            // Skipping map iab in query parameter serialization
            self.us_privacy
                .as_ref()
                .map(|us_privacy| ["us_privacy".to_string(), us_privacy.to_string()].join(",")),
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a Regulations value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for Regulations {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub coppa: Vec<bool>,
            pub eu_privacy: Vec<String>,
            pub gdpr: Vec<bool>,
            pub iab: Vec<std::collections::HashMap<String, serde_json::Value>>,
            pub us_privacy: Vec<String>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err(
                        "Missing value while parsing Regulations".to_string(),
                    )
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    #[allow(clippy::redundant_clone)]
                    "coppa" => intermediate_rep.coppa.push(
                        <bool as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "eu_privacy" => intermediate_rep.eu_privacy.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "gdpr" => intermediate_rep.gdpr.push(
                        <bool as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    "iab" => {
                        return std::result::Result::Err(
                            "Parsing a container in this style is not supported in Regulations"
                                .to_string(),
                        )
                    }
                    #[allow(clippy::redundant_clone)]
                    "us_privacy" => intermediate_rep.us_privacy.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing Regulations".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(Regulations {
            coppa: intermediate_rep.coppa.into_iter().next(),
            eu_privacy: intermediate_rep.eu_privacy.into_iter().next(),
            gdpr: intermediate_rep.gdpr.into_iter().next(),
            iab: intermediate_rep.iab.into_iter().next(),
            us_privacy: intermediate_rep.us_privacy.into_iter().next(),
        })
    }
}

// Methods for converting between header::IntoHeaderValue<Regulations> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<Regulations>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<Regulations>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for Regulations - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<Regulations> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <Regulations as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into Regulations - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<Regulations>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<Regulations>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<Regulations>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<Regulations> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <Regulations as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into Regulations - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct RewardRequest {
    #[serde(rename = "app")]
    pub app: models::App,

    #[serde(rename = "device")]
    pub device: models::Device,

    #[serde(rename = "ext")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ext: Option<String>,

    #[serde(rename = "geo")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub geo: Option<models::Geo>,

    #[serde(rename = "regs")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub regs: Option<models::Regulations>,

    #[serde(rename = "segment")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub segment: Option<models::Segment>,

    #[serde(rename = "session")]
    pub session: models::Session,

    #[serde(rename = "token")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub token: Option<String>,

    #[serde(rename = "user")]
    pub user: models::User,

    #[serde(rename = "bid")]
    #[serde(deserialize_with = "swagger::nullable_format::deserialize_optional_nullable")]
    #[serde(default = "swagger::nullable_format::default_optional_nullable")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub bid: Option<swagger::Nullable<models::Bid>>,

    #[serde(rename = "show")]
    #[serde(deserialize_with = "swagger::nullable_format::deserialize_optional_nullable")]
    #[serde(default = "swagger::nullable_format::default_optional_nullable")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub show: Option<swagger::Nullable<models::Bid>>,
}

impl RewardRequest {
    #[allow(clippy::new_without_default)]
    pub fn new(
        app: models::App,
        device: models::Device,
        session: models::Session,
        user: models::User,
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

/// Converts the RewardRequest value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for RewardRequest {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            // Skipping non-primitive type app in query parameter serialization
            // Skipping non-primitive type device in query parameter serialization
            self.ext
                .as_ref()
                .map(|ext| ["ext".to_string(), ext.to_string()].join(",")),
            // Skipping non-primitive type geo in query parameter serialization
            // Skipping non-primitive type regs in query parameter serialization
            // Skipping non-primitive type segment in query parameter serialization
            // Skipping non-primitive type session in query parameter serialization
            self.token
                .as_ref()
                .map(|token| ["token".to_string(), token.to_string()].join(",")),
            // Skipping non-primitive type user in query parameter serialization
            // Skipping non-primitive type bid in query parameter serialization
            // Skipping non-primitive type show in query parameter serialization
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a RewardRequest value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for RewardRequest {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub app: Vec<models::App>,
            pub device: Vec<models::Device>,
            pub ext: Vec<String>,
            pub geo: Vec<models::Geo>,
            pub regs: Vec<models::Regulations>,
            pub segment: Vec<models::Segment>,
            pub session: Vec<models::Session>,
            pub token: Vec<String>,
            pub user: Vec<models::User>,
            pub bid: Vec<swagger::Nullable<models::Bid>>,
            pub show: Vec<swagger::Nullable<models::Bid>>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err(
                        "Missing value while parsing RewardRequest".to_string(),
                    )
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    #[allow(clippy::redundant_clone)]
                    "app" => intermediate_rep.app.push(
                        <models::App as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "device" => intermediate_rep.device.push(
                        <models::Device as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "ext" => intermediate_rep.ext.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "geo" => intermediate_rep.geo.push(
                        <models::Geo as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "regs" => intermediate_rep.regs.push(
                        <models::Regulations as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "segment" => intermediate_rep.segment.push(
                        <models::Segment as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "session" => intermediate_rep.session.push(
                        <models::Session as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "token" => intermediate_rep.token.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "user" => intermediate_rep.user.push(
                        <models::User as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    "bid" => return std::result::Result::Err(
                        "Parsing a nullable type in this style is not supported in RewardRequest"
                            .to_string(),
                    ),
                    "show" => return std::result::Result::Err(
                        "Parsing a nullable type in this style is not supported in RewardRequest"
                            .to_string(),
                    ),
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing RewardRequest".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(RewardRequest {
            app: intermediate_rep
                .app
                .into_iter()
                .next()
                .ok_or_else(|| "app missing in RewardRequest".to_string())?,
            device: intermediate_rep
                .device
                .into_iter()
                .next()
                .ok_or_else(|| "device missing in RewardRequest".to_string())?,
            ext: intermediate_rep.ext.into_iter().next(),
            geo: intermediate_rep.geo.into_iter().next(),
            regs: intermediate_rep.regs.into_iter().next(),
            segment: intermediate_rep.segment.into_iter().next(),
            session: intermediate_rep
                .session
                .into_iter()
                .next()
                .ok_or_else(|| "session missing in RewardRequest".to_string())?,
            token: intermediate_rep.token.into_iter().next(),
            user: intermediate_rep
                .user
                .into_iter()
                .next()
                .ok_or_else(|| "user missing in RewardRequest".to_string())?,
            bid: std::result::Result::Err(
                "Nullable types not supported in RewardRequest".to_string(),
            )?,
            show: std::result::Result::Err(
                "Nullable types not supported in RewardRequest".to_string(),
            )?,
        })
    }
}

// Methods for converting between header::IntoHeaderValue<RewardRequest> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<RewardRequest>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<RewardRequest>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for RewardRequest - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<RewardRequest> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <RewardRequest as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into RewardRequest - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<RewardRequest>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<RewardRequest>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<RewardRequest>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<RewardRequest> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <RewardRequest as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into RewardRequest - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
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

/// Converts the Segment value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for Segment {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            self.ext
                .as_ref()
                .map(|ext| ["ext".to_string(), ext.to_string()].join(",")),
            self.id
                .as_ref()
                .map(|id| ["id".to_string(), id.to_string()].join(",")),
            self.uid
                .as_ref()
                .map(|uid| ["uid".to_string(), uid.to_string()].join(",")),
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a Segment value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for Segment {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub ext: Vec<String>,
            pub id: Vec<String>,
            pub uid: Vec<String>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err(
                        "Missing value while parsing Segment".to_string(),
                    )
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    #[allow(clippy::redundant_clone)]
                    "ext" => intermediate_rep.ext.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "id" => intermediate_rep.id.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "uid" => intermediate_rep.uid.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing Segment".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(Segment {
            ext: intermediate_rep.ext.into_iter().next(),
            id: intermediate_rep.id.into_iter().next(),
            uid: intermediate_rep.uid.into_iter().next(),
        })
    }
}

// Methods for converting between header::IntoHeaderValue<Segment> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<Segment>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<Segment>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for Segment - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<Segment> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <Segment as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into Segment - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<Segment>>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<Segment>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<Vec<Segment>> {
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<Segment> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <Segment as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into Segment - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
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
    pub launch_monotonic_ts: i32,

    /// Timestamp of the session launch
    #[serde(rename = "launch_ts")]
    pub launch_ts: i32,

    /// Monotonic timestamps when memory warnings occurred
    #[serde(rename = "memory_warnings_monotonic_ts")]
    pub memory_warnings_monotonic_ts: Vec<i32>,

    /// Timestamps when memory warnings occurred
    #[serde(rename = "memory_warnings_ts")]
    pub memory_warnings_ts: Vec<i32>,

    /// Current monotonic timestamp of the session
    #[serde(rename = "monotonic_ts")]
    pub monotonic_ts: i32,

    /// Total size of RAM
    #[serde(rename = "ram_size")]
    pub ram_size: i32,

    /// Amount of RAM used
    #[serde(rename = "ram_used")]
    pub ram_used: i32,

    /// Monotonic timestamp of the session start
    #[serde(rename = "start_monotonic_ts")]
    pub start_monotonic_ts: i32,

    /// Timestamp of the session start
    #[serde(rename = "start_ts")]
    pub start_ts: i32,

    /// Free storage space available
    #[serde(rename = "storage_free")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub storage_free: Option<i32>,

    /// Used storage space
    #[serde(rename = "storage_used")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub storage_used: Option<i32>,

    /// Current timestamp of the session
    #[serde(rename = "ts")]
    pub ts: i32,
}

impl Session {
    #[allow(clippy::new_without_default)]
    pub fn new(
        battery: f64,
        cpu_usage: f64,
        id: uuid::Uuid,
        launch_monotonic_ts: i32,
        launch_ts: i32,
        memory_warnings_monotonic_ts: Vec<i32>,
        memory_warnings_ts: Vec<i32>,
        monotonic_ts: i32,
        ram_size: i32,
        ram_used: i32,
        start_monotonic_ts: i32,
        start_ts: i32,
        ts: i32,
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
            storage_free: None,
            storage_used: None,
            ts,
        }
    }
}

/// Converts the Session value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for Session {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            Some("battery".to_string()),
            Some(self.battery.to_string()),
            Some("cpu_usage".to_string()),
            Some(self.cpu_usage.to_string()),
            // Skipping non-primitive type id in query parameter serialization
            Some("launch_monotonic_ts".to_string()),
            Some(self.launch_monotonic_ts.to_string()),
            Some("launch_ts".to_string()),
            Some(self.launch_ts.to_string()),
            Some("memory_warnings_monotonic_ts".to_string()),
            Some(
                self.memory_warnings_monotonic_ts
                    .iter()
                    .map(|x| x.to_string())
                    .collect::<Vec<_>>()
                    .join(","),
            ),
            Some("memory_warnings_ts".to_string()),
            Some(
                self.memory_warnings_ts
                    .iter()
                    .map(|x| x.to_string())
                    .collect::<Vec<_>>()
                    .join(","),
            ),
            Some("monotonic_ts".to_string()),
            Some(self.monotonic_ts.to_string()),
            Some("ram_size".to_string()),
            Some(self.ram_size.to_string()),
            Some("ram_used".to_string()),
            Some(self.ram_used.to_string()),
            Some("start_monotonic_ts".to_string()),
            Some(self.start_monotonic_ts.to_string()),
            Some("start_ts".to_string()),
            Some(self.start_ts.to_string()),
            self.storage_free.as_ref().map(|storage_free| {
                ["storage_free".to_string(), storage_free.to_string()].join(",")
            }),
            self.storage_used.as_ref().map(|storage_used| {
                ["storage_used".to_string(), storage_used.to_string()].join(",")
            }),
            Some("ts".to_string()),
            Some(self.ts.to_string()),
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a Session value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for Session {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub battery: Vec<f64>,
            pub cpu_usage: Vec<f64>,
            pub id: Vec<uuid::Uuid>,
            pub launch_monotonic_ts: Vec<i32>,
            pub launch_ts: Vec<i32>,
            pub memory_warnings_monotonic_ts: Vec<Vec<i32>>,
            pub memory_warnings_ts: Vec<Vec<i32>>,
            pub monotonic_ts: Vec<i32>,
            pub ram_size: Vec<i32>,
            pub ram_used: Vec<i32>,
            pub start_monotonic_ts: Vec<i32>,
            pub start_ts: Vec<i32>,
            pub storage_free: Vec<i32>,
            pub storage_used: Vec<i32>,
            pub ts: Vec<i32>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err(
                        "Missing value while parsing Session".to_string(),
                    )
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    #[allow(clippy::redundant_clone)]
                    "battery" => intermediate_rep.battery.push(
                        <f64 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "cpu_usage" => intermediate_rep.cpu_usage.push(
                        <f64 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "id" => intermediate_rep.id.push(
                        <uuid::Uuid as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "launch_monotonic_ts" => intermediate_rep.launch_monotonic_ts.push(
                        <i32 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "launch_ts" => intermediate_rep.launch_ts.push(
                        <i32 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    "memory_warnings_monotonic_ts" => {
                        return std::result::Result::Err(
                            "Parsing a container in this style is not supported in Session"
                                .to_string(),
                        )
                    }
                    "memory_warnings_ts" => {
                        return std::result::Result::Err(
                            "Parsing a container in this style is not supported in Session"
                                .to_string(),
                        )
                    }
                    #[allow(clippy::redundant_clone)]
                    "monotonic_ts" => intermediate_rep.monotonic_ts.push(
                        <i32 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "ram_size" => intermediate_rep.ram_size.push(
                        <i32 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "ram_used" => intermediate_rep.ram_used.push(
                        <i32 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "start_monotonic_ts" => intermediate_rep.start_monotonic_ts.push(
                        <i32 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "start_ts" => intermediate_rep.start_ts.push(
                        <i32 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "storage_free" => intermediate_rep.storage_free.push(
                        <i32 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "storage_used" => intermediate_rep.storage_used.push(
                        <i32 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "ts" => intermediate_rep.ts.push(
                        <i32 as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing Session".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(Session {
            battery: intermediate_rep
                .battery
                .into_iter()
                .next()
                .ok_or_else(|| "battery missing in Session".to_string())?,
            cpu_usage: intermediate_rep
                .cpu_usage
                .into_iter()
                .next()
                .ok_or_else(|| "cpu_usage missing in Session".to_string())?,
            id: intermediate_rep
                .id
                .into_iter()
                .next()
                .ok_or_else(|| "id missing in Session".to_string())?,
            launch_monotonic_ts: intermediate_rep
                .launch_monotonic_ts
                .into_iter()
                .next()
                .ok_or_else(|| "launch_monotonic_ts missing in Session".to_string())?,
            launch_ts: intermediate_rep
                .launch_ts
                .into_iter()
                .next()
                .ok_or_else(|| "launch_ts missing in Session".to_string())?,
            memory_warnings_monotonic_ts: intermediate_rep
                .memory_warnings_monotonic_ts
                .into_iter()
                .next()
                .ok_or_else(|| "memory_warnings_monotonic_ts missing in Session".to_string())?,
            memory_warnings_ts: intermediate_rep
                .memory_warnings_ts
                .into_iter()
                .next()
                .ok_or_else(|| "memory_warnings_ts missing in Session".to_string())?,
            monotonic_ts: intermediate_rep
                .monotonic_ts
                .into_iter()
                .next()
                .ok_or_else(|| "monotonic_ts missing in Session".to_string())?,
            ram_size: intermediate_rep
                .ram_size
                .into_iter()
                .next()
                .ok_or_else(|| "ram_size missing in Session".to_string())?,
            ram_used: intermediate_rep
                .ram_used
                .into_iter()
                .next()
                .ok_or_else(|| "ram_used missing in Session".to_string())?,
            start_monotonic_ts: intermediate_rep
                .start_monotonic_ts
                .into_iter()
                .next()
                .ok_or_else(|| "start_monotonic_ts missing in Session".to_string())?,
            start_ts: intermediate_rep
                .start_ts
                .into_iter()
                .next()
                .ok_or_else(|| "start_ts missing in Session".to_string())?,
            storage_free: intermediate_rep.storage_free.into_iter().next(),
            storage_used: intermediate_rep.storage_used.into_iter().next(),
            ts: intermediate_rep
                .ts
                .into_iter()
                .next()
                .ok_or_else(|| "ts missing in Session".to_string())?,
        })
    }
}

// Methods for converting between header::IntoHeaderValue<Session> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<Session>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<Session>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for Session - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<Session> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <Session as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into Session - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<Session>>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<Session>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<Vec<Session>> {
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<Session> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <Session as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into Session - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct ShowRequest {
    #[serde(rename = "app")]
    pub app: models::App,

    #[serde(rename = "device")]
    pub device: models::Device,

    #[serde(rename = "ext")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ext: Option<String>,

    #[serde(rename = "geo")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub geo: Option<models::Geo>,

    #[serde(rename = "regs")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub regs: Option<models::Regulations>,

    #[serde(rename = "segment")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub segment: Option<models::Segment>,

    #[serde(rename = "session")]
    pub session: models::Session,

    #[serde(rename = "token")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub token: Option<String>,

    #[serde(rename = "user")]
    pub user: models::User,

    #[serde(rename = "bid")]
    #[serde(deserialize_with = "swagger::nullable_format::deserialize_optional_nullable")]
    #[serde(default = "swagger::nullable_format::default_optional_nullable")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub bid: Option<swagger::Nullable<models::Bid>>,

    #[serde(rename = "show")]
    #[serde(deserialize_with = "swagger::nullable_format::deserialize_optional_nullable")]
    #[serde(default = "swagger::nullable_format::default_optional_nullable")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub show: Option<swagger::Nullable<models::Bid>>,
}

impl ShowRequest {
    #[allow(clippy::new_without_default)]
    pub fn new(
        app: models::App,
        device: models::Device,
        session: models::Session,
        user: models::User,
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

/// Converts the ShowRequest value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for ShowRequest {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            // Skipping non-primitive type app in query parameter serialization
            // Skipping non-primitive type device in query parameter serialization
            self.ext
                .as_ref()
                .map(|ext| ["ext".to_string(), ext.to_string()].join(",")),
            // Skipping non-primitive type geo in query parameter serialization
            // Skipping non-primitive type regs in query parameter serialization
            // Skipping non-primitive type segment in query parameter serialization
            // Skipping non-primitive type session in query parameter serialization
            self.token
                .as_ref()
                .map(|token| ["token".to_string(), token.to_string()].join(",")),
            // Skipping non-primitive type user in query parameter serialization
            // Skipping non-primitive type bid in query parameter serialization
            // Skipping non-primitive type show in query parameter serialization
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a ShowRequest value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for ShowRequest {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub app: Vec<models::App>,
            pub device: Vec<models::Device>,
            pub ext: Vec<String>,
            pub geo: Vec<models::Geo>,
            pub regs: Vec<models::Regulations>,
            pub segment: Vec<models::Segment>,
            pub session: Vec<models::Session>,
            pub token: Vec<String>,
            pub user: Vec<models::User>,
            pub bid: Vec<swagger::Nullable<models::Bid>>,
            pub show: Vec<swagger::Nullable<models::Bid>>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err(
                        "Missing value while parsing ShowRequest".to_string(),
                    )
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    #[allow(clippy::redundant_clone)]
                    "app" => intermediate_rep.app.push(
                        <models::App as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "device" => intermediate_rep.device.push(
                        <models::Device as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "ext" => intermediate_rep.ext.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "geo" => intermediate_rep.geo.push(
                        <models::Geo as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "regs" => intermediate_rep.regs.push(
                        <models::Regulations as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "segment" => intermediate_rep.segment.push(
                        <models::Segment as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "session" => intermediate_rep.session.push(
                        <models::Session as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "token" => intermediate_rep.token.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "user" => intermediate_rep.user.push(
                        <models::User as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    "bid" => {
                        return std::result::Result::Err(
                            "Parsing a nullable type in this style is not supported in ShowRequest"
                                .to_string(),
                        )
                    }
                    "show" => {
                        return std::result::Result::Err(
                            "Parsing a nullable type in this style is not supported in ShowRequest"
                                .to_string(),
                        )
                    }
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing ShowRequest".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(ShowRequest {
            app: intermediate_rep
                .app
                .into_iter()
                .next()
                .ok_or_else(|| "app missing in ShowRequest".to_string())?,
            device: intermediate_rep
                .device
                .into_iter()
                .next()
                .ok_or_else(|| "device missing in ShowRequest".to_string())?,
            ext: intermediate_rep.ext.into_iter().next(),
            geo: intermediate_rep.geo.into_iter().next(),
            regs: intermediate_rep.regs.into_iter().next(),
            segment: intermediate_rep.segment.into_iter().next(),
            session: intermediate_rep
                .session
                .into_iter()
                .next()
                .ok_or_else(|| "session missing in ShowRequest".to_string())?,
            token: intermediate_rep.token.into_iter().next(),
            user: intermediate_rep
                .user
                .into_iter()
                .next()
                .ok_or_else(|| "user missing in ShowRequest".to_string())?,
            bid: std::result::Result::Err(
                "Nullable types not supported in ShowRequest".to_string(),
            )?,
            show: std::result::Result::Err(
                "Nullable types not supported in ShowRequest".to_string(),
            )?,
        })
    }
}

// Methods for converting between header::IntoHeaderValue<ShowRequest> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<ShowRequest>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<ShowRequest>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for ShowRequest - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<ShowRequest> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <ShowRequest as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into ShowRequest - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<ShowRequest>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<ShowRequest>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<ShowRequest>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<ShowRequest> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <ShowRequest as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into ShowRequest - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize)]
pub struct Stats();

/// Converts the Stats value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl ::std::string::ToString for Stats {
    fn to_string(&self) -> String {
        // ToString for this model is not supported
        "".to_string()
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a Stats value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl ::std::str::FromStr for Stats {
    type Err = &'static str;

    fn from_str(_s: &str) -> std::result::Result<Self, Self::Err> {
        std::result::Result::Err("Parsing Stats is not supported")
    }
}

// Methods for converting between header::IntoHeaderValue<Stats> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<Stats>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<Stats>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for Stats - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<Stats> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => match <Stats as std::str::FromStr>::from_str(value) {
                std::result::Result::Ok(value) => {
                    std::result::Result::Ok(header::IntoHeaderValue(value))
                }
                std::result::Result::Err(err) => std::result::Result::Err(format!(
                    "Unable to convert header value '{}' into Stats - {}",
                    value, err
                )),
            },
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<Stats>>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<Stats>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<Vec<Stats>> {
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<Stats> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <Stats as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into Stats - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct StatsRequest {
    #[serde(rename = "app")]
    pub app: models::App,

    #[serde(rename = "device")]
    pub device: models::Device,

    #[serde(rename = "ext")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ext: Option<String>,

    #[serde(rename = "geo")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub geo: Option<models::Geo>,

    #[serde(rename = "regs")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub regs: Option<models::Regulations>,

    #[serde(rename = "segment")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub segment: Option<models::Segment>,

    #[serde(rename = "session")]
    pub session: models::Session,

    #[serde(rename = "token")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub token: Option<String>,

    #[serde(rename = "user")]
    pub user: models::User,

    #[serde(rename = "stats")]
    pub stats: swagger::Nullable<models::Stats>,
}

impl StatsRequest {
    #[allow(clippy::new_without_default)]
    pub fn new(
        app: models::App,
        device: models::Device,
        session: models::Session,
        user: models::User,
        stats: swagger::Nullable<models::Stats>,
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

/// Converts the StatsRequest value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for StatsRequest {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            // Skipping non-primitive type app in query parameter serialization
            // Skipping non-primitive type device in query parameter serialization
            self.ext
                .as_ref()
                .map(|ext| ["ext".to_string(), ext.to_string()].join(",")),
            // Skipping non-primitive type geo in query parameter serialization
            // Skipping non-primitive type regs in query parameter serialization
            // Skipping non-primitive type segment in query parameter serialization
            // Skipping non-primitive type session in query parameter serialization
            self.token
                .as_ref()
                .map(|token| ["token".to_string(), token.to_string()].join(",")),
            // Skipping non-primitive type user in query parameter serialization
            // Skipping non-primitive type stats in query parameter serialization
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a StatsRequest value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for StatsRequest {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub app: Vec<models::App>,
            pub device: Vec<models::Device>,
            pub ext: Vec<String>,
            pub geo: Vec<models::Geo>,
            pub regs: Vec<models::Regulations>,
            pub segment: Vec<models::Segment>,
            pub session: Vec<models::Session>,
            pub token: Vec<String>,
            pub user: Vec<models::User>,
            pub stats: Vec<swagger::Nullable<models::Stats>>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err(
                        "Missing value while parsing StatsRequest".to_string(),
                    )
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    #[allow(clippy::redundant_clone)]
                    "app" => intermediate_rep.app.push(
                        <models::App as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "device" => intermediate_rep.device.push(
                        <models::Device as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "ext" => intermediate_rep.ext.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "geo" => intermediate_rep.geo.push(
                        <models::Geo as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "regs" => intermediate_rep.regs.push(
                        <models::Regulations as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "segment" => intermediate_rep.segment.push(
                        <models::Segment as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "session" => intermediate_rep.session.push(
                        <models::Session as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "token" => intermediate_rep.token.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "user" => intermediate_rep.user.push(
                        <models::User as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    "stats" => return std::result::Result::Err(
                        "Parsing a nullable type in this style is not supported in StatsRequest"
                            .to_string(),
                    ),
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing StatsRequest".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(StatsRequest {
            app: intermediate_rep
                .app
                .into_iter()
                .next()
                .ok_or_else(|| "app missing in StatsRequest".to_string())?,
            device: intermediate_rep
                .device
                .into_iter()
                .next()
                .ok_or_else(|| "device missing in StatsRequest".to_string())?,
            ext: intermediate_rep.ext.into_iter().next(),
            geo: intermediate_rep.geo.into_iter().next(),
            regs: intermediate_rep.regs.into_iter().next(),
            segment: intermediate_rep.segment.into_iter().next(),
            session: intermediate_rep
                .session
                .into_iter()
                .next()
                .ok_or_else(|| "session missing in StatsRequest".to_string())?,
            token: intermediate_rep.token.into_iter().next(),
            user: intermediate_rep
                .user
                .into_iter()
                .next()
                .ok_or_else(|| "user missing in StatsRequest".to_string())?,
            stats: std::result::Result::Err(
                "Nullable types not supported in StatsRequest".to_string(),
            )?,
        })
    }
}

// Methods for converting between header::IntoHeaderValue<StatsRequest> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<StatsRequest>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<StatsRequest>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for StatsRequest - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<StatsRequest> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <StatsRequest as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into StatsRequest - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<StatsRequest>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<StatsRequest>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<StatsRequest>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<StatsRequest> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <StatsRequest as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into StatsRequest - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
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

/// Converts the SuccessResponse value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for SuccessResponse {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> =
            vec![Some("success".to_string()), Some(self.success.to_string())];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a SuccessResponse value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for SuccessResponse {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub success: Vec<bool>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err(
                        "Missing value while parsing SuccessResponse".to_string(),
                    )
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    #[allow(clippy::redundant_clone)]
                    "success" => intermediate_rep.success.push(
                        <bool as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing SuccessResponse".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(SuccessResponse {
            success: intermediate_rep
                .success
                .into_iter()
                .next()
                .ok_or_else(|| "success missing in SuccessResponse".to_string())?,
        })
    }
}

// Methods for converting between header::IntoHeaderValue<SuccessResponse> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<SuccessResponse>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<SuccessResponse>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for SuccessResponse - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<SuccessResponse>
{
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <SuccessResponse as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into SuccessResponse - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<SuccessResponse>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<SuccessResponse>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<SuccessResponse>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<SuccessResponse> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <SuccessResponse as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into SuccessResponse - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
        }
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

/// Converts the User value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for User {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            // Skipping map consent in query parameter serialization
            self.coppa
                .as_ref()
                .map(|coppa| ["coppa".to_string(), coppa.to_string()].join(",")),
            // Skipping non-primitive type idfa in query parameter serialization
            // Skipping non-primitive type idfv in query parameter serialization
            // Skipping non-primitive type idg in query parameter serialization
            Some("tracking_authorization_status".to_string()),
            Some(self.tracking_authorization_status.to_string()),
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a User value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for User {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub consent: Vec<std::collections::HashMap<String, serde_json::Value>>,
            pub coppa: Vec<bool>,
            pub idfa: Vec<uuid::Uuid>,
            pub idfv: Vec<uuid::Uuid>,
            pub idg: Vec<uuid::Uuid>,
            pub tracking_authorization_status: Vec<String>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err("Missing value while parsing User".to_string())
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    "consent" => {
                        return std::result::Result::Err(
                            "Parsing a container in this style is not supported in User"
                                .to_string(),
                        )
                    }
                    #[allow(clippy::redundant_clone)]
                    "coppa" => intermediate_rep.coppa.push(
                        <bool as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "idfa" => intermediate_rep.idfa.push(
                        <uuid::Uuid as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "idfv" => intermediate_rep.idfv.push(
                        <uuid::Uuid as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "idg" => intermediate_rep.idg.push(
                        <uuid::Uuid as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "tracking_authorization_status" => {
                        intermediate_rep.tracking_authorization_status.push(
                            <String as std::str::FromStr>::from_str(val)
                                .map_err(|x| x.to_string())?,
                        )
                    }
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing User".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(User {
            consent: intermediate_rep.consent.into_iter().next(),
            coppa: intermediate_rep.coppa.into_iter().next(),
            idfa: intermediate_rep.idfa.into_iter().next(),
            idfv: intermediate_rep.idfv.into_iter().next(),
            idg: intermediate_rep.idg.into_iter().next(),
            tracking_authorization_status: intermediate_rep
                .tracking_authorization_status
                .into_iter()
                .next()
                .ok_or_else(|| "tracking_authorization_status missing in User".to_string())?,
        })
    }
}

// Methods for converting between header::IntoHeaderValue<User> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<User>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<User>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for User - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<User> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => match <User as std::str::FromStr>::from_str(value) {
                std::result::Result::Ok(value) => {
                    std::result::Result::Ok(header::IntoHeaderValue(value))
                }
                std::result::Result::Err(err) => std::result::Result::Err(format!(
                    "Unable to convert header value '{}' into User - {}",
                    value, err
                )),
            },
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<User>>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<User>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<Vec<User>> {
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<User> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <User as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into User - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
        }
    }
}

#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize, validator::Validate)]
pub struct WinRequest {
    #[serde(rename = "app")]
    pub app: models::App,

    #[serde(rename = "device")]
    pub device: models::Device,

    #[serde(rename = "ext")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ext: Option<String>,

    #[serde(rename = "geo")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub geo: Option<models::Geo>,

    #[serde(rename = "regs")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub regs: Option<models::Regulations>,

    #[serde(rename = "segment")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub segment: Option<models::Segment>,

    #[serde(rename = "session")]
    pub session: models::Session,

    #[serde(rename = "token")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub token: Option<String>,

    #[serde(rename = "user")]
    pub user: models::User,

    #[serde(rename = "bid")]
    #[serde(deserialize_with = "swagger::nullable_format::deserialize_optional_nullable")]
    #[serde(default = "swagger::nullable_format::default_optional_nullable")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub bid: Option<swagger::Nullable<models::Bid>>,

    #[serde(rename = "show")]
    #[serde(deserialize_with = "swagger::nullable_format::deserialize_optional_nullable")]
    #[serde(default = "swagger::nullable_format::default_optional_nullable")]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub show: Option<swagger::Nullable<models::Bid>>,
}

impl WinRequest {
    #[allow(clippy::new_without_default)]
    pub fn new(
        app: models::App,
        device: models::Device,
        session: models::Session,
        user: models::User,
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

/// Converts the WinRequest value to the Query Parameters representation (style=form, explode=false)
/// specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde serializer
impl std::string::ToString for WinRequest {
    fn to_string(&self) -> String {
        let params: Vec<Option<String>> = vec![
            // Skipping non-primitive type app in query parameter serialization
            // Skipping non-primitive type device in query parameter serialization
            self.ext
                .as_ref()
                .map(|ext| ["ext".to_string(), ext.to_string()].join(",")),
            // Skipping non-primitive type geo in query parameter serialization
            // Skipping non-primitive type regs in query parameter serialization
            // Skipping non-primitive type segment in query parameter serialization
            // Skipping non-primitive type session in query parameter serialization
            self.token
                .as_ref()
                .map(|token| ["token".to_string(), token.to_string()].join(",")),
            // Skipping non-primitive type user in query parameter serialization
            // Skipping non-primitive type bid in query parameter serialization
            // Skipping non-primitive type show in query parameter serialization
        ];

        params.into_iter().flatten().collect::<Vec<_>>().join(",")
    }
}

/// Converts Query Parameters representation (style=form, explode=false) to a WinRequest value
/// as specified in https://swagger.io/docs/specification/serialization/
/// Should be implemented in a serde deserializer
impl std::str::FromStr for WinRequest {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        /// An intermediate representation of the struct to use for parsing.
        #[derive(Default)]
        #[allow(dead_code)]
        struct IntermediateRep {
            pub app: Vec<models::App>,
            pub device: Vec<models::Device>,
            pub ext: Vec<String>,
            pub geo: Vec<models::Geo>,
            pub regs: Vec<models::Regulations>,
            pub segment: Vec<models::Segment>,
            pub session: Vec<models::Session>,
            pub token: Vec<String>,
            pub user: Vec<models::User>,
            pub bid: Vec<swagger::Nullable<models::Bid>>,
            pub show: Vec<swagger::Nullable<models::Bid>>,
        }

        let mut intermediate_rep = IntermediateRep::default();

        // Parse into intermediate representation
        let mut string_iter = s.split(',');
        let mut key_result = string_iter.next();

        while key_result.is_some() {
            let val = match string_iter.next() {
                Some(x) => x,
                None => {
                    return std::result::Result::Err(
                        "Missing value while parsing WinRequest".to_string(),
                    )
                }
            };

            if let Some(key) = key_result {
                #[allow(clippy::match_single_binding)]
                match key {
                    #[allow(clippy::redundant_clone)]
                    "app" => intermediate_rep.app.push(
                        <models::App as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "device" => intermediate_rep.device.push(
                        <models::Device as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "ext" => intermediate_rep.ext.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "geo" => intermediate_rep.geo.push(
                        <models::Geo as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "regs" => intermediate_rep.regs.push(
                        <models::Regulations as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "segment" => intermediate_rep.segment.push(
                        <models::Segment as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "session" => intermediate_rep.session.push(
                        <models::Session as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "token" => intermediate_rep.token.push(
                        <String as std::str::FromStr>::from_str(val).map_err(|x| x.to_string())?,
                    ),
                    #[allow(clippy::redundant_clone)]
                    "user" => intermediate_rep.user.push(
                        <models::User as std::str::FromStr>::from_str(val)
                            .map_err(|x| x.to_string())?,
                    ),
                    "bid" => {
                        return std::result::Result::Err(
                            "Parsing a nullable type in this style is not supported in WinRequest"
                                .to_string(),
                        )
                    }
                    "show" => {
                        return std::result::Result::Err(
                            "Parsing a nullable type in this style is not supported in WinRequest"
                                .to_string(),
                        )
                    }
                    _ => {
                        return std::result::Result::Err(
                            "Unexpected key while parsing WinRequest".to_string(),
                        )
                    }
                }
            }

            // Get the next key
            key_result = string_iter.next();
        }

        // Use the intermediate representation to return the struct
        std::result::Result::Ok(WinRequest {
            app: intermediate_rep
                .app
                .into_iter()
                .next()
                .ok_or_else(|| "app missing in WinRequest".to_string())?,
            device: intermediate_rep
                .device
                .into_iter()
                .next()
                .ok_or_else(|| "device missing in WinRequest".to_string())?,
            ext: intermediate_rep.ext.into_iter().next(),
            geo: intermediate_rep.geo.into_iter().next(),
            regs: intermediate_rep.regs.into_iter().next(),
            segment: intermediate_rep.segment.into_iter().next(),
            session: intermediate_rep
                .session
                .into_iter()
                .next()
                .ok_or_else(|| "session missing in WinRequest".to_string())?,
            token: intermediate_rep.token.into_iter().next(),
            user: intermediate_rep
                .user
                .into_iter()
                .next()
                .ok_or_else(|| "user missing in WinRequest".to_string())?,
            bid: std::result::Result::Err(
                "Nullable types not supported in WinRequest".to_string(),
            )?,
            show: std::result::Result::Err(
                "Nullable types not supported in WinRequest".to_string(),
            )?,
        })
    }
}

// Methods for converting between header::IntoHeaderValue<WinRequest> and hyper::header::HeaderValue

impl std::convert::TryFrom<header::IntoHeaderValue<WinRequest>> for hyper::header::HeaderValue {
    type Error = String;

    fn try_from(
        hdr_value: header::IntoHeaderValue<WinRequest>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_value = hdr_value.to_string();
        match hyper::header::HeaderValue::from_str(&hdr_value) {
            std::result::Result::Ok(value) => std::result::Result::Ok(value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Invalid header value for WinRequest - value: {} is invalid {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue> for header::IntoHeaderValue<WinRequest> {
    type Error = String;

    fn try_from(hdr_value: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_value.to_str() {
            std::result::Result::Ok(value) => {
                match <WinRequest as std::str::FromStr>::from_str(value) {
                    std::result::Result::Ok(value) => {
                        std::result::Result::Ok(header::IntoHeaderValue(value))
                    }
                    std::result::Result::Err(err) => std::result::Result::Err(format!(
                        "Unable to convert header value '{}' into WinRequest - {}",
                        value, err
                    )),
                }
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert header: {:?} to string: {}",
                hdr_value, e
            )),
        }
    }
}

impl std::convert::TryFrom<header::IntoHeaderValue<Vec<WinRequest>>>
    for hyper::header::HeaderValue
{
    type Error = String;

    fn try_from(
        hdr_values: header::IntoHeaderValue<Vec<WinRequest>>,
    ) -> std::result::Result<Self, Self::Error> {
        let hdr_values: Vec<String> = hdr_values
            .0
            .into_iter()
            .map(|hdr_value| hdr_value.to_string())
            .collect();

        match hyper::header::HeaderValue::from_str(&hdr_values.join(", ")) {
            std::result::Result::Ok(hdr_value) => std::result::Result::Ok(hdr_value),
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to convert {:?} into a header - {}",
                hdr_values, e
            )),
        }
    }
}

impl std::convert::TryFrom<hyper::header::HeaderValue>
    for header::IntoHeaderValue<Vec<WinRequest>>
{
    type Error = String;

    fn try_from(hdr_values: hyper::header::HeaderValue) -> std::result::Result<Self, Self::Error> {
        match hdr_values.to_str() {
            std::result::Result::Ok(hdr_values) => {
                let hdr_values: std::vec::Vec<WinRequest> = hdr_values
                    .split(',')
                    .filter_map(|hdr_value| match hdr_value.trim() {
                        "" => std::option::Option::None,
                        hdr_value => std::option::Option::Some({
                            match <WinRequest as std::str::FromStr>::from_str(hdr_value) {
                                std::result::Result::Ok(value) => std::result::Result::Ok(value),
                                std::result::Result::Err(err) => std::result::Result::Err(format!(
                                    "Unable to convert header value '{}' into WinRequest - {}",
                                    hdr_value, err
                                )),
                            }
                        }),
                    })
                    .collect::<std::result::Result<std::vec::Vec<_>, String>>()?;

                std::result::Result::Ok(header::IntoHeaderValue(hdr_values))
            }
            std::result::Result::Err(e) => std::result::Result::Err(format!(
                "Unable to parse header: {:?} as a string - {}",
                hdr_values, e
            )),
        }
    }
}
