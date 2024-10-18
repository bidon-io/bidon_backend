use crate::{models, ServiceError};
use swagger::{ApiError, ContextWrapper};
use async_trait::async_trait;
use std::error::Error;
use std::task::{Context, Poll};
use serde::{Deserialize, Serialize};

#[derive(Debug, PartialEq, Serialize, Deserialize)]
pub enum GetConfigResponse {
    /// Config response
    ConfigResponse
    (models::ConfigResponse)
}

#[derive(Debug, PartialEq, Serialize, Deserialize)]
#[must_use]
pub enum GetOpenApiSpecResponse {
    /// OpenAPI JSON specification
    OpenAPIJSONSpecification(serde_json::Value),
    /// Error
    Error(models::Error),
}

#[derive(Debug, PartialEq, Serialize, Deserialize)]
pub enum GetAuctionResponse {
    /// Auction response
    AuctionResponse(models::AuctionResponse)
}

#[derive(Debug, PartialEq, Serialize, Deserialize)]
pub enum PostClickResponse {
    /// Click response
    ClickResponse(models::SuccessResponse)
}

#[derive(Debug, PartialEq, Serialize, Deserialize)]
pub enum PostLossResponse {
    /// Loss response
    LossResponse(models::SuccessResponse)
}

#[derive(Debug, PartialEq, Serialize, Deserialize)]
pub enum PostRewardResponse {
    /// Reward response
    RewardResponse(models::SuccessResponse)
}

#[derive(Debug, PartialEq, Serialize, Deserialize)]
pub enum PostShowResponse {
    /// Show response
    ShowResponse(models::SuccessResponse)
}

#[derive(Debug, PartialEq, Serialize, Deserialize)]
pub enum PostStatsResponse {
    /// Stats response
    StatsResponse(models::SuccessResponse)
}

#[derive(Debug, PartialEq, Serialize, Deserialize)]
pub enum PostWinResponse {
    /// Win response
    WinResponse(models::SuccessResponse)
}

/// API
#[async_trait]
#[allow(clippy::too_many_arguments, clippy::ptr_arg)]
pub trait Api<C: Send + Sync> {
    fn poll_ready(&self, _cx: &mut Context) -> Poll<Result<(), Box<dyn Error + Send + Sync + 'static>>> {
        Poll::Ready(Ok(()))
    }

    /// Get config
    async fn get_config(&self, config_request: models::ConfigRequest, context: &C)
                        -> Result<GetConfigResponse, ApiError>;

    /// Get OpenAPI specification
    async fn get_open_api_spec(&self, context: &C) -> Result<GetOpenApiSpecResponse, ApiError>;

    /// Auction
    async fn get_auction(&self, ad_type: models::GetAuctionAdTypeParameter,
                         auction_request: models::AuctionRequest,
                         context: &C) -> Result<GetAuctionResponse, ApiError>;

    /// Click
    async fn post_click(
        &self,
        ad_type: models::GetAuctionAdTypeParameter,
        click_request: models::ClickRequest,
        context: &C) -> Result<PostClickResponse, ApiError>;

    /// Loss
    async fn post_loss(
        &self,
        ad_type: models::GetAuctionAdTypeParameter,
        loss_request: models::LossRequest,
        context: &C) -> Result<PostLossResponse, ApiError>;

    /// Reward
    async fn post_reward(
        &self,
        ad_type: models::PostRewardAdTypeParameter,
        reward_request: models::RewardRequest,
        context: &C) -> Result<PostRewardResponse, ApiError>;

    /// Show
    async fn post_show(
        &self,
        ad_type: models::GetAuctionAdTypeParameter,
        show_request: models::ShowRequest,
        context: &C) -> Result<PostShowResponse, ApiError>;

    /// Stats
    async fn post_stats(
        &self,
        ad_type: models::GetAuctionAdTypeParameter,
        stats_request: models::StatsRequest,
        context: &C) -> Result<PostStatsResponse, ApiError>;

    /// Win
    async fn post_win(
        &self,
        ad_type: models::GetAuctionAdTypeParameter,
        win_request: models::WinRequest,
        context: &C) -> Result<PostWinResponse, ApiError>;
}
