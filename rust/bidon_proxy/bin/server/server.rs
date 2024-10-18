//! Main library entry point for bidon implementation.

#![allow(unused_imports)]

use async_trait::async_trait;
use futures::{future, Stream, StreamExt, TryFutureExt, TryStreamExt};
use hyper::server::conn::Http;
use hyper::service::Service;
use log::info;
use std::future::Future;
use std::marker::PhantomData;
use std::net::SocketAddr;
use std::sync::{Arc, Mutex};
use std::task::{Context, Poll};
use swagger::auth::MakeAllowAllAuthenticator;
use swagger::{Has, XSpanIdString};
// use tokio::net::TcpListener;
use std::net::TcpListener;

#[cfg(not(any(target_os = "macos", target_os = "windows", target_os = "ios")))]
use openssl::ssl::{Ssl, SslAcceptor, SslAcceptorBuilder, SslFiletype, SslMethod};

use bidon::models;
use bidon::context::MyEmpContext;
use bidon::bidon_version::XBidonVersionString;

/// Builds an SSL implementation for Simple HTTPS from some hard-coded file names
pub async fn create(addr: &str, https: bool) {
    let addr = addr.parse().expect("Failed to parse bind address");

    let server = Server::new();

    let service = MakeService::new(server);

    let service = MakeAllowAllAuthenticator::new(service, "cosmo");

    #[allow(unused_mut)]
    let mut service = bidon::server::context::MakeAddContext::<_, MyEmpContext>::new(service);

    if https {
        #[cfg(any(target_os = "macos", target_os = "windows", target_os = "ios"))]
        {
            unimplemented!("SSL is not implemented for the examples on MacOS, Windows or iOS");
        }

        #[cfg(not(any(target_os = "macos", target_os = "windows", target_os = "ios")))]
        {
            let mut ssl = SslAcceptor::mozilla_intermediate_v5(SslMethod::tls())
                .expect("Failed to create SSL Acceptor");

            // Server authentication
            ssl.set_private_key_file("examples/server-key.pem", SslFiletype::PEM)
                .expect("Failed to set private key");
            ssl.set_certificate_chain_file("examples/server-chain.pem")
                .expect("Failed to set certificate chain");
            ssl.check_private_key()
                .expect("Failed to check private key");

            let tls_acceptor = ssl.build();
            let tcp_listener = TcpListener::bind(&addr).await.unwrap();

            info!("Starting a server (with https)");
            loop {
                if let Ok((tcp, _)) = tcp_listener.accept().await {
                    let ssl = Ssl::new(tls_acceptor.context()).unwrap();
                    let addr = tcp.peer_addr().expect("Unable to get remote address");
                    let service = service.call(addr);

                    tokio::spawn(async move {
                        let tls = tokio_openssl::SslStream::new(ssl, tcp).map_err(|_| ())?;
                        let service = service.await.map_err(|_| ())?;

                        Http::new()
                            .serve_connection(tls, service)
                            .await
                            .map_err(|_| ())
                    });
                }
            }
        }
    } else {
        info!("Starting a server (over http, so no TLS)");
        // Using HTTP
        hyper::server::Server::bind(&addr)
            .serve(service)
            .await
            .unwrap()
    }
}

#[derive(Copy, Clone)]
pub struct Server<C> {
    marker: PhantomData<C>,
}

impl<C> Server<C> {
    pub fn new() -> Self {
        Server {
            marker: PhantomData,
        }
    }
}

use crate::server_auth;
use jsonwebtoken::{
    decode, encode, errors::Error as JwtError, Algorithm, DecodingKey, EncodingKey, Header,
    TokenData, Validation,
};
use serde::{Deserialize, Serialize};
use swagger::auth::Authorization;

use bidon::server::MakeService;
use bidon::{
    Api, GetAuctionResponse, GetConfigResponse, GetOpenApiSpecResponse, PostClickResponse,
    PostLossResponse, PostRewardResponse, PostShowResponse, PostStatsResponse, PostWinResponse,
};
use std::error::Error;
use swagger::ApiError;

#[async_trait]
impl<C> Api<C> for Server<C>
where
    C: Has<XSpanIdString> + Has<XBidonVersionString> + Send + Sync,
{
    /// Get config
    async fn get_config(
        &self,
        config_request: models::ConfigRequest,
        context: &C,
    ) -> Result<GetConfigResponse, ApiError> {
        info!(
            "get_config(\"{}\", {:?}) - X-Span-ID: {:?}",
            <C as Has<XBidonVersionString>>::get(context).0.clone(),
            config_request,
            <C as Has<XSpanIdString>>::get(context).0.clone()
        );
        Err(ApiError("Api-Error: Operation is NOT implemented".into()))
    }

    /// Get OpenAPI specification
    async fn get_open_api_spec(&self, context: &C) -> Result<GetOpenApiSpecResponse, ApiError> {
        info!(
            "get_open_api_spec() - X-Span-ID: {:?}",
            <C as Has<XSpanIdString>>::get(context).0.clone()
        );
        Err(ApiError("Api-Error: Operation is NOT implemented".into()))
    }

    /// Auction
    async fn get_auction(
        &self,
        ad_type: models::GetAuctionAdTypeParameter,
        auction_request: models::AuctionRequest,
        context: &C,
    ) -> Result<GetAuctionResponse, ApiError> {
        info!(
            "get_auction(\"{}\", {:?}, {:?}) - X-Span-ID: {:?}",
            <C as Has<XBidonVersionString>>::get(context).0.clone(),
            ad_type,
            auction_request,
            <C as Has<XSpanIdString>>::get(context).0.clone()
        );
        Err(ApiError("Api-Error: Operation is NOT implemented".into()))
    }

    /// Click
    async fn post_click(
        &self,
        ad_type: models::GetAuctionAdTypeParameter,
        click_request: models::ClickRequest,
        context: &C,
    ) -> Result<PostClickResponse, ApiError> {
        info!(
            "post_click(\"{}\", {:?}, {:?}) - X-Span-ID: {:?}",
            <C as Has<XBidonVersionString>>::get(context).0.clone(),
            ad_type,
            click_request,
            <C as Has<XSpanIdString>>::get(context).0.clone()

        );
        Err(ApiError("Api-Error: Operation is NOT implemented".into()))
    }

    /// Loss
    async fn post_loss(
        &self,
        ad_type: models::GetAuctionAdTypeParameter,
        loss_request: models::LossRequest,
        context: &C,
    ) -> Result<PostLossResponse, ApiError> {
        info!(
            "post_loss(\"{}\", {:?}, {:?}) - X-Span-ID: {:?}",
            <C as Has<XBidonVersionString>>::get(context).0.clone(),
            ad_type,
            loss_request,
            <C as Has<XSpanIdString>>::get(context).0.clone()
        );
        Err(ApiError("Api-Error: Operation is NOT implemented".into()))
    }

    /// Reward
    async fn post_reward(
        &self,
        ad_type: models::PostRewardAdTypeParameter,
        reward_request: models::RewardRequest,
        context: &C,
    ) -> Result<PostRewardResponse, ApiError> {
        info!(
            "post_reward(\"{}\", {:?}, {:?}) - X-Span-ID: {:?}",
            <C as Has<XBidonVersionString>>::get(context).0.clone(),
            ad_type,
            reward_request,
            <C as Has<XSpanIdString>>::get(context).0.clone()
        );
        Err(ApiError("Api-Error: Operation is NOT implemented".into()))
    }

    /// Show
    async fn post_show(
        &self,
        ad_type: models::GetAuctionAdTypeParameter,
        show_request: models::ShowRequest,
        context: &C,
    ) -> Result<PostShowResponse, ApiError> {
        info!(
            "post_show(\"{}\", {:?}, {:?}) - X-Span-ID: {:?}",
            <C as Has<XBidonVersionString>>::get(context).0.clone(),
            ad_type,
            show_request,
            <C as Has<XSpanIdString>>::get(context).0.clone()
        );
        Err(ApiError("Api-Error: Operation is NOT implemented".into()))
    }

    /// Stats
    async fn post_stats(
        &self,
        ad_type: models::GetAuctionAdTypeParameter,
        stats_request: models::StatsRequest,
        context: &C,
    ) -> Result<PostStatsResponse, ApiError> {
        info!(
            "post_stats(\"{}\", {:?}, {:?}) - X-Span-ID: {:?}",
            <C as Has<XBidonVersionString>>::get(context).0.clone(),
            ad_type,
            stats_request,
            <C as Has<XSpanIdString>>::get(context).0.clone()
        );
        Err(ApiError("Api-Error: Operation is NOT implemented".into()))
    }

    /// Win
    async fn post_win(
        &self,
        ad_type: models::GetAuctionAdTypeParameter,
        win_request: models::WinRequest,
        context: &C,
    ) -> Result<PostWinResponse, ApiError> {
        info!(
            "post_win(\"{}\", {:?}, {:?}) - X-Span-ID: {:?}",
            <C as Has<XBidonVersionString>>::get(context).0.clone(),
            ad_type,
            win_request,
            <C as Has<XSpanIdString>>::get(context).0.clone()
        );
        Err(ApiError("Api-Error: Operation is NOT implemented".into()))
    }
}
