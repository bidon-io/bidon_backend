use crate::bidding::error::BiddingError;
use crate::bidding::BiddingService;
use crate::com::iabtechlab::openrtb::v3::Openrtb;
use crate::org::bidon::proto::v1::bidding_service_client::BiddingServiceClient;
use tonic::transport::Channel;
use tonic::Request;

#[derive(Debug, Clone)]
pub struct ProxyBiddingService {
    grpc_client: BiddingServiceClient<Channel>,
}

impl ProxyBiddingService {
    pub async fn new(grpc_url: &'static str) -> Result<Self, Box<dyn std::error::Error>> {
        let grpc_client = BiddingServiceClient::connect(grpc_url).await?;
        Ok(ProxyBiddingService { grpc_client })
    }
}

#[async_trait::async_trait]
impl BiddingService for ProxyBiddingService {
    async fn bid(&self, request: Openrtb) -> Result<Openrtb, BiddingError> {
        let grpc_request = Request::new(request);
        let grpc_response = self
            .grpc_client
            .clone() // Cloning is required here, because Tonic gRPC clients are mutable. Cloning is cheap.
            .bid(grpc_request)
            .await
            .map_err(BiddingError::GrpcError)?;

        Ok(grpc_response.into_inner())
    }
}
