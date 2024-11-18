#![allow(
    missing_docs,
    trivial_casts,
    unused_variables,
    unused_mut,
    unused_imports,
    unused_extern_crates,
    unused_attributes,
    non_camel_case_types
)]
#![allow(clippy::derive_partial_eq_without_eq, clippy::disallowed_names)]

use async_trait::async_trait;
use futures::Stream;
use serde::{Deserialize, Serialize};
use std::error::Error;
use std::task::{Context, Poll};
use swagger::{ApiError, ContextWrapper};

pub const BASE_PATH: &str = "";
pub const API_VERSION: &str = "1.0.0";

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

pub mod galaxy {
    pub mod v1 {
        tonic::include_proto!("galaxy.v1");
        pub mod bidon {
            tonic::include_proto!("galaxy.v1.bidon");
        }
        pub mod context {
            tonic::include_proto!("galaxy.v1.context");
        }
    }
}

pub mod models;

pub(crate) mod header;

pub mod auction;

pub mod controllers {
    pub mod adapter;
    pub mod auction;
}

pub mod bidon_version;
