use crate::bidding::Api;
use crate::bidding::BiddingError;
use galaxy_bidon::com::iabtechlab::openrtb::v3::Openrtb;
use galaxy_bidon::galaxy::v1::bidding_service_client::BiddingServiceClient;
use tonic::transport::Channel;
use tonic::Request;

pub struct ProxyBiddingService {
    grpc_client: BiddingServiceClient<Channel>,
}

impl ProxyBiddingService {
    pub async fn new(grpc_url: String) -> Result<Self, Box<dyn std::error::Error>> {
        let grpc_client = BiddingServiceClient::connect(grpc_url).await?;
        Ok(ProxyBiddingService { grpc_client })
    }
}

#[async_trait::async_trait]
impl Api for ProxyBiddingService {
    async fn bid(&mut self, bidding_request: Openrtb) -> Result<Openrtb, BiddingError> {
        let request = Request::new(bidding_request);
        let response = self
            .grpc_client
            .bid(request)
            .await
            .map_err(|e| BiddingError::new(e.to_string()))?;
        Ok(response.into_inner())
    }
}
