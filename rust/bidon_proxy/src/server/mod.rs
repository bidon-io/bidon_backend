use futures::{future, future::BoxFuture, future::FutureExt, stream, stream::TryStreamExt, Stream};
use hyper::{Body, HeaderMap, Request, Response, StatusCode};
use hyper::header::{HeaderName, HeaderValue, CONTENT_TYPE};
use log::warn;
#[allow(unused_imports)]
use std::convert::{TryFrom, TryInto};
use std::error::Error;
use std::future::Future;
use std::marker::PhantomData;
use std::task::{Context, Poll};
use swagger::{ApiError, BodyExt, Has, RequestParser, XSpanIdString};
pub use swagger::auth::Authorization;
use swagger::auth::Scopes;
use url::form_urlencoded;

#[allow(unused_imports)]
use crate::{header, models, AuthenticationApi};

pub use crate::context;

type ServiceFuture = BoxFuture<'static, Result<Response<Body>, crate::ServiceError>>;

use crate::{Api,
            GetAuctionResponse,
            GetConfigResponse,
            GetOpenApiSpecResponse,
            PostClickResponse,
            PostLossResponse,
            PostRewardResponse,
            PostShowResponse,
            PostStatsResponse,
            PostWinResponse
};
use crate::bidon_version::XBidonVersionString;

mod server_auth;

mod paths {
    use lazy_static::lazy_static;

    lazy_static! {
        pub static ref GLOBAL_REGEX_SET: regex::RegexSet = regex::RegexSet::new(vec![
            r"^/openapi.json$",
            r"^/v2/auction/(?P<ad_type>[^/?#]*)$",
            r"^/v2/click/(?P<ad_type>[^/?#]*)$",
            r"^/v2/config$",
            r"^/v2/loss/(?P<ad_type>[^/?#]*)$",
            r"^/v2/reward/(?P<ad_type>[^/?#]*)$",
            r"^/v2/show/(?P<ad_type>[^/?#]*)$",
            r"^/v2/stats/(?P<ad_type>[^/?#]*)$",
            r"^/v2/win/(?P<ad_type>[^/?#]*)$"
        ])
        .expect("Unable to create global regex set");
    }
    pub(crate) static ID_OPENAPI_JSON: usize = 0;
    pub(crate) static ID_V2_AUCTION_AD_TYPE: usize = 1;
    lazy_static! {
        pub static ref REGEX_V2_AUCTION_AD_TYPE: regex::Regex =
            #[allow(clippy::invalid_regex)]
            regex::Regex::new(r"^/v2/auction/(?P<ad_type>[^/?#]*)$")
                .expect("Unable to create regex for V2_AUCTION_AD_TYPE");
    }
    pub(crate) static ID_V2_CLICK_AD_TYPE: usize = 2;
    lazy_static! {
        pub static ref REGEX_V2_CLICK_AD_TYPE: regex::Regex =
            #[allow(clippy::invalid_regex)]
            regex::Regex::new(r"^/v2/click/(?P<ad_type>[^/?#]*)$")
                .expect("Unable to create regex for V2_CLICK_AD_TYPE");
    }
    pub(crate) static ID_V2_CONFIG: usize = 3;
    pub(crate) static ID_V2_LOSS_AD_TYPE: usize = 4;
    lazy_static! {
        pub static ref REGEX_V2_LOSS_AD_TYPE: regex::Regex =
            #[allow(clippy::invalid_regex)]
            regex::Regex::new(r"^/v2/loss/(?P<ad_type>[^/?#]*)$")
                .expect("Unable to create regex for V2_LOSS_AD_TYPE");
    }
    pub(crate) static ID_V2_REWARD_AD_TYPE: usize = 5;
    lazy_static! {
        pub static ref REGEX_V2_REWARD_AD_TYPE: regex::Regex =
            #[allow(clippy::invalid_regex)]
            regex::Regex::new(r"^/v2/reward/(?P<ad_type>[^/?#]*)$")
                .expect("Unable to create regex for V2_REWARD_AD_TYPE");
    }
    pub(crate) static ID_V2_SHOW_AD_TYPE: usize = 6;
    lazy_static! {
        pub static ref REGEX_V2_SHOW_AD_TYPE: regex::Regex =
            #[allow(clippy::invalid_regex)]
            regex::Regex::new(r"^/v2/show/(?P<ad_type>[^/?#]*)$")
                .expect("Unable to create regex for V2_SHOW_AD_TYPE");
    }
    pub(crate) static ID_V2_STATS_AD_TYPE: usize = 7;
    lazy_static! {
        pub static ref REGEX_V2_STATS_AD_TYPE: regex::Regex =
            #[allow(clippy::invalid_regex)]
            regex::Regex::new(r"^/v2/stats/(?P<ad_type>[^/?#]*)$")
                .expect("Unable to create regex for V2_STATS_AD_TYPE");
    }
    pub(crate) static ID_V2_WIN_AD_TYPE: usize = 8;
    lazy_static! {
        pub static ref REGEX_V2_WIN_AD_TYPE: regex::Regex =
            #[allow(clippy::invalid_regex)]
            regex::Regex::new(r"^/v2/win/(?P<ad_type>[^/?#]*)$")
                .expect("Unable to create regex for V2_WIN_AD_TYPE");
    }
}


pub struct MakeService<T, C> where
    T: Api<C> + Clone + Send + 'static,
    C: Has<XSpanIdString> + Has<XBidonVersionString> + Send + Sync + 'static
{
    api_impl: T,
    marker: PhantomData<C>,
}

impl<T, C> MakeService<T, C> where
    T: Api<C> + Clone + Send + 'static,
    C: Has<XSpanIdString> + Has<XBidonVersionString>  + Send + Sync + 'static
{
    pub fn new(api_impl: T) -> Self {
        MakeService {
            api_impl,
            marker: PhantomData
        }
    }
}

impl<T, C, Target> hyper::service::Service<Target> for MakeService<T, C> where
    T: Api<C> + Clone + Send + 'static,
    C: Has<XSpanIdString> + Has<XBidonVersionString> + Send + Sync + 'static
{
    type Response = Service<T, C>;
    type Error = crate::ServiceError;
    type Future = future::Ready<Result<Self::Response, Self::Error>>;

    fn poll_ready(&mut self, cx: &mut Context<'_>) -> Poll<Result<(), Self::Error>> {
        Poll::Ready(Ok(()))
    }

    fn call(&mut self, target: Target) -> Self::Future {
        let service = Service::new(self.api_impl.clone());

        future::ok(service)
    }
}

fn method_not_allowed() -> Result<Response<Body>, crate::ServiceError> {
    Ok(
        Response::builder().status(StatusCode::METHOD_NOT_ALLOWED)
            .body(Body::empty())
            .expect("Unable to create Method Not Allowed response")
    )
}

pub struct Service<T, C> where
    T: Api<C> + Clone + Send + 'static,
    C: Has<XSpanIdString> + Has<XBidonVersionString> + Send + Sync + 'static
{
    api_impl: T,
    marker: PhantomData<C>,
}

impl<T, C> Service<T, C> where
    T: Api<C> + Clone + Send + 'static,
    C: Has<XSpanIdString> + Has<XBidonVersionString> + Send + Sync + 'static
{
    pub fn new(api_impl: T) -> Self {
        Service {
            api_impl,
            marker: PhantomData
        }
    }
}

impl<T, C> Clone for Service<T, C> where
    T: Api<C> + Clone + Send + 'static,
    C: Has<XSpanIdString> + Has<XBidonVersionString>  + Send + Sync + 'static
{
    fn clone(&self) -> Self {
        Service {
            api_impl: self.api_impl.clone(),
            marker: self.marker,
        }
    }
}

impl<T, C> hyper::service::Service<(Request<Body>, C)> for Service<T, C> where
    T: Api<C> + Clone + Send + Sync + 'static,
    C: Has<XSpanIdString> + Has<XBidonVersionString>  + Send + Sync + 'static
{
    type Response = Response<Body>;
    type Error = crate::ServiceError;
    type Future = ServiceFuture;

    fn poll_ready(&mut self, cx: &mut Context) -> Poll<Result<(), Self::Error>> {
        self.api_impl.poll_ready(cx)
    }

    fn call(&mut self, req: (Request<Body>, C)) -> Self::Future {
        async fn run<T, C>(
            mut api_impl: T,
            req: (Request<Body>, C),
        ) -> Result<Response<Body>, crate::ServiceError> where
            T: Api<C> + Clone + Send + 'static,
            C: Has<XSpanIdString> + Has<XBidonVersionString>  + Send + Sync + 'static
        {
            let (request, context) = req;
            let (parts, body) = request.into_parts();
            let (method, uri, headers) = (parts.method, parts.uri, parts.headers);
            let path = paths::GLOBAL_REGEX_SET.matches(uri.path());

            match method {

            // GetConfig - POST /v2/config
            hyper::Method::POST if path.matched(paths::ID_V2_CONFIG) => {

                // Handle body parameters (note that non-required body parameters will ignore garbage
                // values, rather than causing a 400 response). Produce warning header and logs for
                // any unused fields.
                let result = body.into_raw().await;
                match result {
                     Ok(body) => {
                                let mut unused_elements : Vec<String> = vec![];
                                let param_config_request: Option<models::ConfigRequest> = if !body.is_empty() {
                                    let deserializer = &mut serde_json::Deserializer::from_slice(&*body);
                                    match serde_ignored::deserialize(deserializer, |path| {
                                            warn!("Ignoring unknown field in body: {}", path);
                                            unused_elements.push(path.to_string());
                                    }) {
                                        Ok(param_config_request) => param_config_request,
                                        Err(e) => return Ok(Response::builder()
                                                        .status(StatusCode::BAD_REQUEST)
                                                        .body(Body::from(format!("Couldn't parse body parameter ConfigRequest - doesn't match schema: {}", e)))
                                                        .expect("Unable to create Bad Request response for invalid body parameter ConfigRequest due to schema")),
                                    }
                                } else {
                                    None
                                };
                                let param_config_request = match param_config_request {
                                    Some(param_config_request) => param_config_request,
                                    None => return Ok(Response::builder()
                                                        .status(StatusCode::BAD_REQUEST)
                                                        .body(Body::from("Missing required body parameter ConfigRequest"))
                                                        .expect("Unable to create Bad Request response for missing body parameter ConfigRequest")),
                                };


                                let result = api_impl.get_config(
                                            param_config_request,
                                        &context
                                    ).await;
                                let mut response = Response::new(Body::empty());
                                response.headers_mut().insert(
                                            HeaderName::from_static("x-span-id"),
                                            HeaderValue::from_str((&context as &dyn Has<XSpanIdString>).get().0.clone().as_str())
                                                .expect("Unable to create X-Span-ID header value"));

                                        if !unused_elements.is_empty() {
                                            response.headers_mut().insert(
                                                HeaderName::from_static("warning"),
                                                HeaderValue::from_str(format!("Ignoring unknown fields in body: {:?}", unused_elements).as_str())
                                                    .expect("Unable to create Warning header value"));
                                        }
                                        match result {
                                            Ok(rsp) => match rsp {
                                                GetConfigResponse::ConfigResponse
                                                    (body)
                                                => {
                                                    *response.status_mut() = StatusCode::from_u16(200).expect("Unable to turn 200 into a StatusCode");
                                                    response.headers_mut().insert(
                                                        CONTENT_TYPE,
                                                        HeaderValue::from_str("application/json")
                                                            .expect("Unable to create Content-Type header for application/json"));
                                                    // JSON Body
                                                    let body = serde_json::to_string(&body).expect("impossible to fail to serialize");
                                                    *response.body_mut() = Body::from(body);

                                                },
                                            },
                                            Err(_) => {
                                                // Application code returned an error. This should not happen, as the implementation should
                                                // return a valid response.
                                                *response.status_mut() = StatusCode::INTERNAL_SERVER_ERROR;
                                                *response.body_mut() = Body::from("An internal error occurred");
                                            },
                                        }

                                        Ok(response)
                            },
                            Err(e) => Ok(Response::builder()
                                                .status(StatusCode::BAD_REQUEST)
                                                .body(Body::from(format!("Unable to read body: {}", e)))
                                                .expect("Unable to create Bad Request response due to unable to read body")),
                        }
            },

            // GetOpenApiSpec - GET /openapi.json
            hyper::Method::GET if path.matched(paths::ID_OPENAPI_JSON) => {
                                let result = api_impl.get_open_api_spec(
                                        &context
                                    ).await;
                                let mut response = Response::new(Body::empty());
                                response.headers_mut().insert(
                                            HeaderName::from_static("x-span-id"),
                                            HeaderValue::from_str((&context as &dyn Has<XSpanIdString>).get().0.clone().as_str())
                                                .expect("Unable to create X-Span-ID header value"));

                                        match result {
                                            Ok(rsp) => match rsp {
                                                GetOpenApiSpecResponse::OpenAPIJSONSpecification
                                                    (body)
                                                => {
                                                    *response.status_mut() = StatusCode::from_u16(200).expect("Unable to turn 200 into a StatusCode");
                                                    response.headers_mut().insert(
                                                        CONTENT_TYPE,
                                                        HeaderValue::from_str("application/json")
                                                            .expect("Unable to create Content-Type header for application/json"));
                                                    // JSON Body
                                                    let body = serde_json::to_string(&body).expect("impossible to fail to serialize");
                                                    *response.body_mut() = Body::from(body);

                                                },
                                                GetOpenApiSpecResponse::Error
                                                    (body)
                                                => {
                                                    *response.status_mut() = StatusCode::from_u16(0).expect("Unable to turn 0 into a StatusCode");
                                                    response.headers_mut().insert(
                                                        CONTENT_TYPE,
                                                        HeaderValue::from_str("application/problem+json")
                                                            .expect("Unable to create Content-Type header for application/problem+json"));
                                                    // JSON Body
                                                    let body = serde_json::to_string(&body).expect("impossible to fail to serialize");
                                                    *response.body_mut() = Body::from(body);

                                                },
                                            },
                                            Err(_) => {
                                                // Application code returned an error. This should not happen, as the implementation should
                                                // return a valid response.
                                                *response.status_mut() = StatusCode::INTERNAL_SERVER_ERROR;
                                                *response.body_mut() = Body::from("An internal error occurred");
                                            },
                                        }

                                        Ok(response)
            },

            // GetAuction - POST /v2/auction/{ad_type}
            hyper::Method::POST if path.matched(paths::ID_V2_AUCTION_AD_TYPE) => {
                // Path parameters
                let path: &str = uri.path();
                let path_params =
                    paths::REGEX_V2_AUCTION_AD_TYPE
                    .captures(path)
                    .unwrap_or_else(||
                        panic!("Path {} matched RE V2_AUCTION_AD_TYPE in set but failed match against \"{}\"", path, paths::REGEX_V2_AUCTION_AD_TYPE.as_str())
                    );

                let param_ad_type = match percent_encoding::percent_decode(path_params["ad_type"].as_bytes()).decode_utf8() {
                    Ok(param_ad_type) => match param_ad_type.parse::<models::GetAuctionAdTypeParameter>() {
                        Ok(param_ad_type) => param_ad_type,
                        Err(e) => return Ok(Response::builder()
                                        .status(StatusCode::BAD_REQUEST)
                                        .body(Body::from(format!("Couldn't parse path parameter ad_type: {}", e)))
                                        .expect("Unable to create Bad Request response for invalid path parameter")),
                    },
                    Err(_) => return Ok(Response::builder()
                                        .status(StatusCode::BAD_REQUEST)
                                        .body(Body::from(format!("Couldn't percent-decode path parameter as UTF-8: {}", &path_params["ad_type"])))
                                        .expect("Unable to create Bad Request response for invalid percent decode"))
                };

                // Handle body parameters (note that non-required body parameters will ignore garbage
                // values, rather than causing a 400 response). Produce warning header and logs for
                // any unused fields.
                let result = body.into_raw().await;
                match result {
                     Ok(body) => {
                                let mut unused_elements : Vec<String> = vec![];
                                let param_auction_request: Option<models::AuctionRequest> = if !body.is_empty() {
                                    let deserializer = &mut serde_json::Deserializer::from_slice(&*body);
                                    match serde_ignored::deserialize(deserializer, |path| {
                                            warn!("Ignoring unknown field in body: {}", path);
                                            unused_elements.push(path.to_string());
                                    }) {
                                        Ok(param_auction_request) => param_auction_request,
                                        Err(e) => return Ok(Response::builder()
                                                        .status(StatusCode::BAD_REQUEST)
                                                        .body(Body::from(format!("Couldn't parse body parameter AuctionRequest - doesn't match schema: {}", e)))
                                                        .expect("Unable to create Bad Request response for invalid body parameter AuctionRequest due to schema")),
                                    }
                                } else {
                                    None
                                };
                                let param_auction_request = match param_auction_request {
                                    Some(param_auction_request) => param_auction_request,
                                    None => return Ok(Response::builder()
                                                        .status(StatusCode::BAD_REQUEST)
                                                        .body(Body::from("Missing required body parameter AuctionRequest"))
                                                        .expect("Unable to create Bad Request response for missing body parameter AuctionRequest")),
                                };


                                let result = api_impl.get_auction(
                                            param_ad_type,
                                            param_auction_request,
                                        &context
                                    ).await;
                                let mut response = Response::new(Body::empty());
                                response.headers_mut().insert(
                                            HeaderName::from_static("x-span-id"),
                                            HeaderValue::from_str((&context as &dyn Has<XSpanIdString>).get().0.clone().as_str())
                                                .expect("Unable to create X-Span-ID header value"));

                                        if !unused_elements.is_empty() {
                                            response.headers_mut().insert(
                                                HeaderName::from_static("warning"),
                                                HeaderValue::from_str(format!("Ignoring unknown fields in body: {:?}", unused_elements).as_str())
                                                    .expect("Unable to create Warning header value"));
                                        }
                                        match result {
                                            Ok(rsp) => match rsp {
                                                GetAuctionResponse::AuctionResponse
                                                    (body)
                                                => {
                                                    *response.status_mut() = StatusCode::from_u16(200).expect("Unable to turn 200 into a StatusCode");
                                                    response.headers_mut().insert(
                                                        CONTENT_TYPE,
                                                        HeaderValue::from_str("application/json")
                                                            .expect("Unable to create Content-Type header for application/json"));
                                                    // JSON Body
                                                    let body = serde_json::to_string(&body).expect("impossible to fail to serialize");
                                                    *response.body_mut() = Body::from(body);

                                                },
                                            },
                                            Err(_) => {
                                                // Application code returned an error. This should not happen, as the implementation should
                                                // return a valid response.
                                                *response.status_mut() = StatusCode::INTERNAL_SERVER_ERROR;
                                                *response.body_mut() = Body::from("An internal error occurred");
                                            },
                                        }

                                        Ok(response)
                            },
                            Err(e) => Ok(Response::builder()
                                                .status(StatusCode::BAD_REQUEST)
                                                .body(Body::from(format!("Unable to read body: {}", e)))
                                                .expect("Unable to create Bad Request response due to unable to read body")),
                        }
            },

            // PostClick - POST /v2/click/{ad_type}
            hyper::Method::POST if path.matched(paths::ID_V2_CLICK_AD_TYPE) => {
                // Path parameters
                let path: &str = uri.path();
                let path_params =
                    paths::REGEX_V2_CLICK_AD_TYPE
                    .captures(path)
                    .unwrap_or_else(||
                        panic!("Path {} matched RE V2_CLICK_AD_TYPE in set but failed match against \"{}\"", path, paths::REGEX_V2_CLICK_AD_TYPE.as_str())
                    );

                let param_ad_type = match percent_encoding::percent_decode(path_params["ad_type"].as_bytes()).decode_utf8() {
                    Ok(param_ad_type) => match param_ad_type.parse::<models::GetAuctionAdTypeParameter>() {
                        Ok(param_ad_type) => param_ad_type,
                        Err(e) => return Ok(Response::builder()
                                        .status(StatusCode::BAD_REQUEST)
                                        .body(Body::from(format!("Couldn't parse path parameter ad_type: {}", e)))
                                        .expect("Unable to create Bad Request response for invalid path parameter")),
                    },
                    Err(_) => return Ok(Response::builder()
                                        .status(StatusCode::BAD_REQUEST)
                                        .body(Body::from(format!("Couldn't percent-decode path parameter as UTF-8: {}", &path_params["ad_type"])))
                                        .expect("Unable to create Bad Request response for invalid percent decode"))
                };

                // Handle body parameters (note that non-required body parameters will ignore garbage
                // values, rather than causing a 400 response). Produce warning header and logs for
                // any unused fields.
                let result = body.into_raw().await;
                match result {
                     Ok(body) => {
                                let mut unused_elements : Vec<String> = vec![];
                                let param_click_request: Option<models::ClickRequest> = if !body.is_empty() {
                                    let deserializer = &mut serde_json::Deserializer::from_slice(&*body);
                                    match serde_ignored::deserialize(deserializer, |path| {
                                            warn!("Ignoring unknown field in body: {}", path);
                                            unused_elements.push(path.to_string());
                                    }) {
                                        Ok(param_click_request) => param_click_request,
                                        Err(e) => return Ok(Response::builder()
                                                        .status(StatusCode::BAD_REQUEST)
                                                        .body(Body::from(format!("Couldn't parse body parameter ClickRequest - doesn't match schema: {}", e)))
                                                        .expect("Unable to create Bad Request response for invalid body parameter ClickRequest due to schema")),
                                    }
                                } else {
                                    None
                                };
                                let param_click_request = match param_click_request {
                                    Some(param_click_request) => param_click_request,
                                    None => return Ok(Response::builder()
                                                        .status(StatusCode::BAD_REQUEST)
                                                        .body(Body::from("Missing required body parameter ClickRequest"))
                                                        .expect("Unable to create Bad Request response for missing body parameter ClickRequest")),
                                };


                                let result = api_impl.post_click(
                                            param_ad_type,
                                            param_click_request,
                                        &context
                                    ).await;
                                let mut response = Response::new(Body::empty());
                                response.headers_mut().insert(
                                            HeaderName::from_static("x-span-id"),
                                            HeaderValue::from_str((&context as &dyn Has<XSpanIdString>).get().0.clone().as_str())
                                                .expect("Unable to create X-Span-ID header value"));

                                        if !unused_elements.is_empty() {
                                            response.headers_mut().insert(
                                                HeaderName::from_static("warning"),
                                                HeaderValue::from_str(format!("Ignoring unknown fields in body: {:?}", unused_elements).as_str())
                                                    .expect("Unable to create Warning header value"));
                                        }
                                        match result {
                                            Ok(rsp) => match rsp {
                                                PostClickResponse::ClickResponse
                                                    (body)
                                                => {
                                                    *response.status_mut() = StatusCode::from_u16(200).expect("Unable to turn 200 into a StatusCode");
                                                    response.headers_mut().insert(
                                                        CONTENT_TYPE,
                                                        HeaderValue::from_str("application/json")
                                                            .expect("Unable to create Content-Type header for application/json"));
                                                    // JSON Body
                                                    let body = serde_json::to_string(&body).expect("impossible to fail to serialize");
                                                    *response.body_mut() = Body::from(body);

                                                },
                                            },
                                            Err(_) => {
                                                // Application code returned an error. This should not happen, as the implementation should
                                                // return a valid response.
                                                *response.status_mut() = StatusCode::INTERNAL_SERVER_ERROR;
                                                *response.body_mut() = Body::from("An internal error occurred");
                                            },
                                        }

                                        Ok(response)
                            },
                            Err(e) => Ok(Response::builder()
                                                .status(StatusCode::BAD_REQUEST)
                                                .body(Body::from(format!("Unable to read body: {}", e)))
                                                .expect("Unable to create Bad Request response due to unable to read body")),
                        }
            },

            // PostLoss - POST /v2/loss/{ad_type}
            hyper::Method::POST if path.matched(paths::ID_V2_LOSS_AD_TYPE) => {
                // Path parameters
                let path: &str = uri.path();
                let path_params =
                    paths::REGEX_V2_LOSS_AD_TYPE
                    .captures(path)
                    .unwrap_or_else(||
                        panic!("Path {} matched RE V2_LOSS_AD_TYPE in set but failed match against \"{}\"", path, paths::REGEX_V2_LOSS_AD_TYPE.as_str())
                    );

                let param_ad_type = match percent_encoding::percent_decode(path_params["ad_type"].as_bytes()).decode_utf8() {
                    Ok(param_ad_type) => match param_ad_type.parse::<models::GetAuctionAdTypeParameter>() {
                        Ok(param_ad_type) => param_ad_type,
                        Err(e) => return Ok(Response::builder()
                                        .status(StatusCode::BAD_REQUEST)
                                        .body(Body::from(format!("Couldn't parse path parameter ad_type: {}", e)))
                                        .expect("Unable to create Bad Request response for invalid path parameter")),
                    },
                    Err(_) => return Ok(Response::builder()
                                        .status(StatusCode::BAD_REQUEST)
                                        .body(Body::from(format!("Couldn't percent-decode path parameter as UTF-8: {}", &path_params["ad_type"])))
                                        .expect("Unable to create Bad Request response for invalid percent decode"))
                };

                // Handle body parameters (note that non-required body parameters will ignore garbage
                // values, rather than causing a 400 response). Produce warning header and logs for
                // any unused fields.
                let result = body.into_raw().await;
                match result {
                     Ok(body) => {
                                let mut unused_elements : Vec<String> = vec![];
                                let param_loss_request: Option<models::LossRequest> = if !body.is_empty() {
                                    let deserializer = &mut serde_json::Deserializer::from_slice(&*body);
                                    match serde_ignored::deserialize(deserializer, |path| {
                                            warn!("Ignoring unknown field in body: {}", path);
                                            unused_elements.push(path.to_string());
                                    }) {
                                        Ok(param_loss_request) => param_loss_request,
                                        Err(e) => return Ok(Response::builder()
                                                        .status(StatusCode::BAD_REQUEST)
                                                        .body(Body::from(format!("Couldn't parse body parameter LossRequest - doesn't match schema: {}", e)))
                                                        .expect("Unable to create Bad Request response for invalid body parameter LossRequest due to schema")),
                                    }
                                } else {
                                    None
                                };
                                let param_loss_request = match param_loss_request {
                                    Some(param_loss_request) => param_loss_request,
                                    None => return Ok(Response::builder()
                                                        .status(StatusCode::BAD_REQUEST)
                                                        .body(Body::from("Missing required body parameter LossRequest"))
                                                        .expect("Unable to create Bad Request response for missing body parameter LossRequest")),
                                };


                                let result = api_impl.post_loss(
                                            param_ad_type,
                                            param_loss_request,
                                        &context
                                    ).await;
                                let mut response = Response::new(Body::empty());
                                response.headers_mut().insert(
                                            HeaderName::from_static("x-span-id"),
                                            HeaderValue::from_str((&context as &dyn Has<XSpanIdString>).get().0.clone().as_str())
                                                .expect("Unable to create X-Span-ID header value"));

                                        if !unused_elements.is_empty() {
                                            response.headers_mut().insert(
                                                HeaderName::from_static("warning"),
                                                HeaderValue::from_str(format!("Ignoring unknown fields in body: {:?}", unused_elements).as_str())
                                                    .expect("Unable to create Warning header value"));
                                        }
                                        match result {
                                            Ok(rsp) => match rsp {
                                                PostLossResponse::LossResponse
                                                    (body)
                                                => {
                                                    *response.status_mut() = StatusCode::from_u16(200).expect("Unable to turn 200 into a StatusCode");
                                                    response.headers_mut().insert(
                                                        CONTENT_TYPE,
                                                        HeaderValue::from_str("application/json")
                                                            .expect("Unable to create Content-Type header for application/json"));
                                                    // JSON Body
                                                    let body = serde_json::to_string(&body).expect("impossible to fail to serialize");
                                                    *response.body_mut() = Body::from(body);

                                                },
                                            },
                                            Err(_) => {
                                                // Application code returned an error. This should not happen, as the implementation should
                                                // return a valid response.
                                                *response.status_mut() = StatusCode::INTERNAL_SERVER_ERROR;
                                                *response.body_mut() = Body::from("An internal error occurred");
                                            },
                                        }

                                        Ok(response)
                            },
                            Err(e) => Ok(Response::builder()
                                                .status(StatusCode::BAD_REQUEST)
                                                .body(Body::from(format!("Unable to read body: {}", e)))
                                                .expect("Unable to create Bad Request response due to unable to read body")),
                        }
            },

            // PostReward - POST /v2/reward/{ad_type}
            hyper::Method::POST if path.matched(paths::ID_V2_REWARD_AD_TYPE) => {
                // Path parameters
                let path: &str = uri.path();
                let path_params =
                    paths::REGEX_V2_REWARD_AD_TYPE
                    .captures(path)
                    .unwrap_or_else(||
                        panic!("Path {} matched RE V2_REWARD_AD_TYPE in set but failed match against \"{}\"", path, paths::REGEX_V2_REWARD_AD_TYPE.as_str())
                    );

                let param_ad_type = match percent_encoding::percent_decode(path_params["ad_type"].as_bytes()).decode_utf8() {
                    Ok(param_ad_type) => match param_ad_type.parse::<models::PostRewardAdTypeParameter>() {
                        Ok(param_ad_type) => param_ad_type,
                        Err(e) => return Ok(Response::builder()
                                        .status(StatusCode::BAD_REQUEST)
                                        .body(Body::from(format!("Couldn't parse path parameter ad_type: {}", e)))
                                        .expect("Unable to create Bad Request response for invalid path parameter")),
                    },
                    Err(_) => return Ok(Response::builder()
                                        .status(StatusCode::BAD_REQUEST)
                                        .body(Body::from(format!("Couldn't percent-decode path parameter as UTF-8: {}", &path_params["ad_type"])))
                                        .expect("Unable to create Bad Request response for invalid percent decode"))
                };

                // Handle body parameters (note that non-required body parameters will ignore garbage
                // values, rather than causing a 400 response). Produce warning header and logs for
                // any unused fields.
                let result = body.into_raw().await;
                match result {
                     Ok(body) => {
                                let mut unused_elements : Vec<String> = vec![];
                                let param_reward_request: Option<models::RewardRequest> = if !body.is_empty() {
                                    let deserializer = &mut serde_json::Deserializer::from_slice(&*body);
                                    match serde_ignored::deserialize(deserializer, |path| {
                                            warn!("Ignoring unknown field in body: {}", path);
                                            unused_elements.push(path.to_string());
                                    }) {
                                        Ok(param_reward_request) => param_reward_request,
                                        Err(e) => return Ok(Response::builder()
                                                        .status(StatusCode::BAD_REQUEST)
                                                        .body(Body::from(format!("Couldn't parse body parameter RewardRequest - doesn't match schema: {}", e)))
                                                        .expect("Unable to create Bad Request response for invalid body parameter RewardRequest due to schema")),
                                    }
                                } else {
                                    None
                                };
                                let param_reward_request = match param_reward_request {
                                    Some(param_reward_request) => param_reward_request,
                                    None => return Ok(Response::builder()
                                                        .status(StatusCode::BAD_REQUEST)
                                                        .body(Body::from("Missing required body parameter RewardRequest"))
                                                        .expect("Unable to create Bad Request response for missing body parameter RewardRequest")),
                                };


                                let result = api_impl.post_reward(
                                            param_ad_type,
                                            param_reward_request,
                                        &context
                                    ).await;
                                let mut response = Response::new(Body::empty());
                                response.headers_mut().insert(
                                            HeaderName::from_static("x-span-id"),
                                            HeaderValue::from_str((&context as &dyn Has<XSpanIdString>).get().0.clone().as_str())
                                                .expect("Unable to create X-Span-ID header value"));

                                        if !unused_elements.is_empty() {
                                            response.headers_mut().insert(
                                                HeaderName::from_static("warning"),
                                                HeaderValue::from_str(format!("Ignoring unknown fields in body: {:?}", unused_elements).as_str())
                                                    .expect("Unable to create Warning header value"));
                                        }
                                        match result {
                                            Ok(rsp) => match rsp {
                                                PostRewardResponse::RewardResponse
                                                    (body)
                                                => {
                                                    *response.status_mut() = StatusCode::from_u16(200).expect("Unable to turn 200 into a StatusCode");
                                                    response.headers_mut().insert(
                                                        CONTENT_TYPE,
                                                        HeaderValue::from_str("application/json")
                                                            .expect("Unable to create Content-Type header for application/json"));
                                                    // JSON Body
                                                    let body = serde_json::to_string(&body).expect("impossible to fail to serialize");
                                                    *response.body_mut() = Body::from(body);

                                                },
                                            },
                                            Err(_) => {
                                                // Application code returned an error. This should not happen, as the implementation should
                                                // return a valid response.
                                                *response.status_mut() = StatusCode::INTERNAL_SERVER_ERROR;
                                                *response.body_mut() = Body::from("An internal error occurred");
                                            },
                                        }

                                        Ok(response)
                            },
                            Err(e) => Ok(Response::builder()
                                                .status(StatusCode::BAD_REQUEST)
                                                .body(Body::from(format!("Unable to read body: {}", e)))
                                                .expect("Unable to create Bad Request response due to unable to read body")),
                        }
            },

            // PostShow - POST /v2/show/{ad_type}
            hyper::Method::POST if path.matched(paths::ID_V2_SHOW_AD_TYPE) => {
                // Path parameters
                let path: &str = uri.path();
                let path_params =
                    paths::REGEX_V2_SHOW_AD_TYPE
                    .captures(path)
                    .unwrap_or_else(||
                        panic!("Path {} matched RE V2_SHOW_AD_TYPE in set but failed match against \"{}\"", path, paths::REGEX_V2_SHOW_AD_TYPE.as_str())
                    );

                let param_ad_type = match percent_encoding::percent_decode(path_params["ad_type"].as_bytes()).decode_utf8() {
                    Ok(param_ad_type) => match param_ad_type.parse::<models::GetAuctionAdTypeParameter>() {
                        Ok(param_ad_type) => param_ad_type,
                        Err(e) => return Ok(Response::builder()
                                        .status(StatusCode::BAD_REQUEST)
                                        .body(Body::from(format!("Couldn't parse path parameter ad_type: {}", e)))
                                        .expect("Unable to create Bad Request response for invalid path parameter")),
                    },
                    Err(_) => return Ok(Response::builder()
                                        .status(StatusCode::BAD_REQUEST)
                                        .body(Body::from(format!("Couldn't percent-decode path parameter as UTF-8: {}", &path_params["ad_type"])))
                                        .expect("Unable to create Bad Request response for invalid percent decode"))
                };

                // Handle body parameters (note that non-required body parameters will ignore garbage
                // values, rather than causing a 400 response). Produce warning header and logs for
                // any unused fields.
                let result = body.into_raw().await;
                match result {
                     Ok(body) => {
                                let mut unused_elements : Vec<String> = vec![];
                                let param_show_request: Option<models::ShowRequest> = if !body.is_empty() {
                                    let deserializer = &mut serde_json::Deserializer::from_slice(&*body);
                                    match serde_ignored::deserialize(deserializer, |path| {
                                            warn!("Ignoring unknown field in body: {}", path);
                                            unused_elements.push(path.to_string());
                                    }) {
                                        Ok(param_show_request) => param_show_request,
                                        Err(e) => return Ok(Response::builder()
                                                        .status(StatusCode::BAD_REQUEST)
                                                        .body(Body::from(format!("Couldn't parse body parameter ShowRequest - doesn't match schema: {}", e)))
                                                        .expect("Unable to create Bad Request response for invalid body parameter ShowRequest due to schema")),
                                    }
                                } else {
                                    None
                                };
                                let param_show_request = match param_show_request {
                                    Some(param_show_request) => param_show_request,
                                    None => return Ok(Response::builder()
                                                        .status(StatusCode::BAD_REQUEST)
                                                        .body(Body::from("Missing required body parameter ShowRequest"))
                                                        .expect("Unable to create Bad Request response for missing body parameter ShowRequest")),
                                };


                                let result = api_impl.post_show(
                                            param_ad_type,
                                            param_show_request,
                                        &context
                                    ).await;
                                let mut response = Response::new(Body::empty());
                                response.headers_mut().insert(
                                            HeaderName::from_static("x-span-id"),
                                            HeaderValue::from_str((&context as &dyn Has<XSpanIdString>).get().0.clone().as_str())
                                                .expect("Unable to create X-Span-ID header value"));

                                        if !unused_elements.is_empty() {
                                            response.headers_mut().insert(
                                                HeaderName::from_static("warning"),
                                                HeaderValue::from_str(format!("Ignoring unknown fields in body: {:?}", unused_elements).as_str())
                                                    .expect("Unable to create Warning header value"));
                                        }
                                        match result {
                                            Ok(rsp) => match rsp {
                                                PostShowResponse::ShowResponse
                                                    (body)
                                                => {
                                                    *response.status_mut() = StatusCode::from_u16(200).expect("Unable to turn 200 into a StatusCode");
                                                    response.headers_mut().insert(
                                                        CONTENT_TYPE,
                                                        HeaderValue::from_str("application/json")
                                                            .expect("Unable to create Content-Type header for application/json"));
                                                    // JSON Body
                                                    let body = serde_json::to_string(&body).expect("impossible to fail to serialize");
                                                    *response.body_mut() = Body::from(body);

                                                },
                                            },
                                            Err(_) => {
                                                // Application code returned an error. This should not happen, as the implementation should
                                                // return a valid response.
                                                *response.status_mut() = StatusCode::INTERNAL_SERVER_ERROR;
                                                *response.body_mut() = Body::from("An internal error occurred");
                                            },
                                        }

                                        Ok(response)
                            },
                            Err(e) => Ok(Response::builder()
                                                .status(StatusCode::BAD_REQUEST)
                                                .body(Body::from(format!("Unable to read body: {}", e)))
                                                .expect("Unable to create Bad Request response due to unable to read body")),
                        }
            },

            // PostStats - POST /v2/stats/{ad_type}
            hyper::Method::POST if path.matched(paths::ID_V2_STATS_AD_TYPE) => {
                // Path parameters
                let path: &str = uri.path();
                let path_params =
                    paths::REGEX_V2_STATS_AD_TYPE
                    .captures(path)
                    .unwrap_or_else(||
                        panic!("Path {} matched RE V2_STATS_AD_TYPE in set but failed match against \"{}\"", path, paths::REGEX_V2_STATS_AD_TYPE.as_str())
                    );

                let param_ad_type = match percent_encoding::percent_decode(path_params["ad_type"].as_bytes()).decode_utf8() {
                    Ok(param_ad_type) => match param_ad_type.parse::<models::GetAuctionAdTypeParameter>() {
                        Ok(param_ad_type) => param_ad_type,
                        Err(e) => return Ok(Response::builder()
                                        .status(StatusCode::BAD_REQUEST)
                                        .body(Body::from(format!("Couldn't parse path parameter ad_type: {}", e)))
                                        .expect("Unable to create Bad Request response for invalid path parameter")),
                    },
                    Err(_) => return Ok(Response::builder()
                                        .status(StatusCode::BAD_REQUEST)
                                        .body(Body::from(format!("Couldn't percent-decode path parameter as UTF-8: {}", &path_params["ad_type"])))
                                        .expect("Unable to create Bad Request response for invalid percent decode"))
                };

                // Handle body parameters (note that non-required body parameters will ignore garbage
                // values, rather than causing a 400 response). Produce warning header and logs for
                // any unused fields.
                let result = body.into_raw().await;
                match result {
                     Ok(body) => {
                                let mut unused_elements : Vec<String> = vec![];
                                let param_stats_request: Option<models::StatsRequest> = if !body.is_empty() {
                                    let deserializer = &mut serde_json::Deserializer::from_slice(&*body);
                                    match serde_ignored::deserialize(deserializer, |path| {
                                            warn!("Ignoring unknown field in body: {}", path);
                                            unused_elements.push(path.to_string());
                                    }) {
                                        Ok(param_stats_request) => param_stats_request,
                                        Err(e) => return Ok(Response::builder()
                                                        .status(StatusCode::BAD_REQUEST)
                                                        .body(Body::from(format!("Couldn't parse body parameter StatsRequest - doesn't match schema: {}", e)))
                                                        .expect("Unable to create Bad Request response for invalid body parameter StatsRequest due to schema")),
                                    }
                                } else {
                                    None
                                };
                                let param_stats_request = match param_stats_request {
                                    Some(param_stats_request) => param_stats_request,
                                    None => return Ok(Response::builder()
                                                        .status(StatusCode::BAD_REQUEST)
                                                        .body(Body::from("Missing required body parameter StatsRequest"))
                                                        .expect("Unable to create Bad Request response for missing body parameter StatsRequest")),
                                };


                                let result = api_impl.post_stats(
                                            param_ad_type,
                                            param_stats_request,
                                        &context
                                    ).await;
                                let mut response = Response::new(Body::empty());
                                response.headers_mut().insert(
                                            HeaderName::from_static("x-span-id"),
                                            HeaderValue::from_str((&context as &dyn Has<XSpanIdString>).get().0.clone().as_str())
                                                .expect("Unable to create X-Span-ID header value"));

                                        if !unused_elements.is_empty() {
                                            response.headers_mut().insert(
                                                HeaderName::from_static("warning"),
                                                HeaderValue::from_str(format!("Ignoring unknown fields in body: {:?}", unused_elements).as_str())
                                                    .expect("Unable to create Warning header value"));
                                        }
                                        match result {
                                            Ok(rsp) => match rsp {
                                                PostStatsResponse::StatsResponse
                                                    (body)
                                                => {
                                                    *response.status_mut() = StatusCode::from_u16(200).expect("Unable to turn 200 into a StatusCode");
                                                    response.headers_mut().insert(
                                                        CONTENT_TYPE,
                                                        HeaderValue::from_str("application/json")
                                                            .expect("Unable to create Content-Type header for application/json"));
                                                    // JSON Body
                                                    let body = serde_json::to_string(&body).expect("impossible to fail to serialize");
                                                    *response.body_mut() = Body::from(body);

                                                },
                                            },
                                            Err(_) => {
                                                // Application code returned an error. This should not happen, as the implementation should
                                                // return a valid response.
                                                *response.status_mut() = StatusCode::INTERNAL_SERVER_ERROR;
                                                *response.body_mut() = Body::from("An internal error occurred");
                                            },
                                        }

                                        Ok(response)
                            },
                            Err(e) => Ok(Response::builder()
                                                .status(StatusCode::BAD_REQUEST)
                                                .body(Body::from(format!("Unable to read body: {}", e)))
                                                .expect("Unable to create Bad Request response due to unable to read body")),
                        }
            },

            // PostWin - POST /v2/win/{ad_type}
            hyper::Method::POST if path.matched(paths::ID_V2_WIN_AD_TYPE) => {
                // Path parameters
                let path: &str = uri.path();
                let path_params =
                    paths::REGEX_V2_WIN_AD_TYPE
                    .captures(path)
                    .unwrap_or_else(||
                        panic!("Path {} matched RE V2_WIN_AD_TYPE in set but failed match against \"{}\"", path, paths::REGEX_V2_WIN_AD_TYPE.as_str())
                    );

                let param_ad_type = match percent_encoding::percent_decode(path_params["ad_type"].as_bytes()).decode_utf8() {
                    Ok(param_ad_type) => match param_ad_type.parse::<models::GetAuctionAdTypeParameter>() {
                        Ok(param_ad_type) => param_ad_type,
                        Err(e) => return Ok(Response::builder()
                                        .status(StatusCode::BAD_REQUEST)
                                        .body(Body::from(format!("Couldn't parse path parameter ad_type: {}", e)))
                                        .expect("Unable to create Bad Request response for invalid path parameter")),
                    },
                    Err(_) => return Ok(Response::builder()
                                        .status(StatusCode::BAD_REQUEST)
                                        .body(Body::from(format!("Couldn't percent-decode path parameter as UTF-8: {}", &path_params["ad_type"])))
                                        .expect("Unable to create Bad Request response for invalid percent decode"))
                };

                // Handle body parameters (note that non-required body parameters will ignore garbage
                // values, rather than causing a 400 response). Produce warning header and logs for
                // any unused fields.
                let result = body.into_raw().await;
                match result {
                     Ok(body) => {
                                let mut unused_elements : Vec<String> = vec![];
                                let param_win_request: Option<models::WinRequest> = if !body.is_empty() {
                                    let deserializer = &mut serde_json::Deserializer::from_slice(&*body);
                                    match serde_ignored::deserialize(deserializer, |path| {
                                            warn!("Ignoring unknown field in body: {}", path);
                                            unused_elements.push(path.to_string());
                                    }) {
                                        Ok(param_win_request) => param_win_request,
                                        Err(e) => return Ok(Response::builder()
                                                        .status(StatusCode::BAD_REQUEST)
                                                        .body(Body::from(format!("Couldn't parse body parameter WinRequest - doesn't match schema: {}", e)))
                                                        .expect("Unable to create Bad Request response for invalid body parameter WinRequest due to schema")),
                                    }
                                } else {
                                    None
                                };
                                let param_win_request = match param_win_request {
                                    Some(param_win_request) => param_win_request,
                                    None => return Ok(Response::builder()
                                                        .status(StatusCode::BAD_REQUEST)
                                                        .body(Body::from("Missing required body parameter WinRequest"))
                                                        .expect("Unable to create Bad Request response for missing body parameter WinRequest")),
                                };


                                let result = api_impl.post_win(
                                            param_ad_type,
                                            param_win_request,
                                        &context
                                    ).await;
                                let mut response = Response::new(Body::empty());
                                response.headers_mut().insert(
                                            HeaderName::from_static("x-span-id"),
                                            HeaderValue::from_str((&context as &dyn Has<XSpanIdString>).get().0.clone().as_str())
                                                .expect("Unable to create X-Span-ID header value"));

                                        if !unused_elements.is_empty() {
                                            response.headers_mut().insert(
                                                HeaderName::from_static("warning"),
                                                HeaderValue::from_str(format!("Ignoring unknown fields in body: {:?}", unused_elements).as_str())
                                                    .expect("Unable to create Warning header value"));
                                        }
                                        match result {
                                            Ok(rsp) => match rsp {
                                                PostWinResponse::WinResponse
                                                    (body)
                                                => {
                                                    *response.status_mut() = StatusCode::from_u16(200).expect("Unable to turn 200 into a StatusCode");
                                                    response.headers_mut().insert(
                                                        CONTENT_TYPE,
                                                        HeaderValue::from_str("application/json")
                                                            .expect("Unable to create Content-Type header for application/json"));
                                                    // JSON Body
                                                    let body = serde_json::to_string(&body).expect("impossible to fail to serialize");
                                                    *response.body_mut() = Body::from(body);

                                                },
                                            },
                                            Err(_) => {
                                                // Application code returned an error. This should not happen, as the implementation should
                                                // return a valid response.
                                                *response.status_mut() = StatusCode::INTERNAL_SERVER_ERROR;
                                                *response.body_mut() = Body::from("An internal error occurred");
                                            },
                                        }

                                        Ok(response)
                            },
                            Err(e) => Ok(Response::builder()
                                                .status(StatusCode::BAD_REQUEST)
                                                .body(Body::from(format!("Unable to read body: {}", e)))
                                                .expect("Unable to create Bad Request response due to unable to read body")),
                        }
            },

            _ if path.matched(paths::ID_OPENAPI_JSON) => method_not_allowed(),
            _ if path.matched(paths::ID_V2_AUCTION_AD_TYPE) => method_not_allowed(),
            _ if path.matched(paths::ID_V2_CLICK_AD_TYPE) => method_not_allowed(),
            _ if path.matched(paths::ID_V2_CONFIG) => method_not_allowed(),
            _ if path.matched(paths::ID_V2_LOSS_AD_TYPE) => method_not_allowed(),
            _ if path.matched(paths::ID_V2_REWARD_AD_TYPE) => method_not_allowed(),
            _ if path.matched(paths::ID_V2_SHOW_AD_TYPE) => method_not_allowed(),
            _ if path.matched(paths::ID_V2_STATS_AD_TYPE) => method_not_allowed(),
            _ if path.matched(paths::ID_V2_WIN_AD_TYPE) => method_not_allowed(),
                _ => Ok(Response::builder().status(StatusCode::NOT_FOUND)
                        .body(Body::empty())
                        .expect("Unable to create Not Found response"))
            }
        }
        Box::pin(run(
            self.api_impl.clone(),
            req,
        ))
    }
}

/// Request parser for `Api`.
pub struct ApiRequestParser;
impl<T> RequestParser<T> for ApiRequestParser {
    fn parse_operation_id(request: &Request<T>) -> Option<&'static str> {
        let path = paths::GLOBAL_REGEX_SET.matches(request.uri().path());
        match *request.method() {
            // GetConfig - POST /v2/config
            hyper::Method::POST if path.matched(paths::ID_V2_CONFIG) => Some("GetConfig"),
            // GetOpenApiSpec - GET /openapi.json
            hyper::Method::GET if path.matched(paths::ID_OPENAPI_JSON) => Some("GetOpenApiSpec"),
            // GetAuction - POST /v2/auction/{ad_type}
            hyper::Method::POST if path.matched(paths::ID_V2_AUCTION_AD_TYPE) => Some("GetAuction"),
            // PostClick - POST /v2/click/{ad_type}
            hyper::Method::POST if path.matched(paths::ID_V2_CLICK_AD_TYPE) => Some("PostClick"),
            // PostLoss - POST /v2/loss/{ad_type}
            hyper::Method::POST if path.matched(paths::ID_V2_LOSS_AD_TYPE) => Some("PostLoss"),
            // PostReward - POST /v2/reward/{ad_type}
            hyper::Method::POST if path.matched(paths::ID_V2_REWARD_AD_TYPE) => Some("PostReward"),
            // PostShow - POST /v2/show/{ad_type}
            hyper::Method::POST if path.matched(paths::ID_V2_SHOW_AD_TYPE) => Some("PostShow"),
            // PostStats - POST /v2/stats/{ad_type}
            hyper::Method::POST if path.matched(paths::ID_V2_STATS_AD_TYPE) => Some("PostStats"),
            // PostWin - POST /v2/win/{ad_type}
            hyper::Method::POST if path.matched(paths::ID_V2_WIN_AD_TYPE) => Some("PostWin"),
            _ => None,
        }
    }
}
