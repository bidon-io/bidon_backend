use crate::bidding::BiddingService;
use axum::routing::post;
use axum::Router;

// Define the routes
pub fn create_app<A>(bidding_service: Box<A>) -> Router
where
    A: Clone + 'static,
    A: BiddingService + Send + Sync,
{
    Router::new()
        .route(
            "/v2/auction/:ad_type",
            post(handlers::auction::get_auction_handler),
        )
        .with_state(bidding_service)
}

// mod main;
pub mod bidding;
pub mod extract;
pub mod sdk;

mod adapter;
mod handlers;

pub mod com {
    pub mod iabtechlab {
        pub mod openrtb {
            pub mod v3 {
                tonic::include_proto!("com.iabtechlab.openrtb.v3");
            }
        }
        pub mod adcom {
            pub mod v1 {
                pub mod context {
                    tonic::include_proto!("com.iabtechlab.adcom.v1.context");
                }
                pub mod enums {
                    tonic::include_proto!("com.iabtechlab.adcom.v1.enums");
                }
                pub mod media {
                    tonic::include_proto!("com.iabtechlab.adcom.v1.media");
                }
                pub mod placement {
                    tonic::include_proto!("com.iabtechlab.adcom.v1.placement");
                }
            }
        }
    }
}

pub mod org {
    pub mod bidon {
        pub mod proto {
            pub mod v1 {
                tonic::include_proto!("org.bidon.proto.v1");
                pub mod mediation {
                    tonic::include_proto!("org.bidon.proto.v1.mediation");
                }
                pub mod context {
                    tonic::include_proto!("org.bidon.proto.v1.context");
                }
            }
        }
    }
}
