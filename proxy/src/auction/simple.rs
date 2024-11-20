use crate::auction::Api;
use crate::auction::AuctionError;
use crate::com::iabtechlab::openrtb::v3::Openrtb;
use crate::galaxy::v1::bidding_service_client::BiddingServiceClient;
use tonic::transport::Channel;
use tonic::Request;

pub struct SimpleAuction {
    grpc_client: BiddingServiceClient<Channel>,
}
//todo rename to proxy/ grpc bidding service
impl SimpleAuction {
    pub async fn new(grpc_url: String) -> Result<Self, Box<dyn std::error::Error>> {
        let grpc_client = BiddingServiceClient::connect(grpc_url).await?;
        Ok(SimpleAuction { grpc_client })
    }
}

#[async_trait::async_trait]
impl Api for SimpleAuction {
    async fn bid(&mut self, auction_request: Openrtb) -> Result<Openrtb, AuctionError> {
        let request = Request::new(auction_request);
        let response = self
            .grpc_client
            .bid(request)
            .await
            .map_err(|e| AuctionError::new(e.to_string()))?;
        Ok(response.into_inner())
    }
}
