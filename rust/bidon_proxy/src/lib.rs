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
use std::error::Error;
use std::task::{Context, Poll};
use swagger::{ApiError, ContextWrapper};
use serde::{Deserialize, Serialize};
type ServiceError = Box<dyn Error + Send + Sync + 'static>;

pub const BASE_PATH: &str = "";
pub const API_VERSION: &str = "1.0.0";

pub mod galaxy {
    pub mod v1 {
        tonic::include_proto!("galaxy.v1");
        pub mod bidon {
            tonic::include_proto!("galaxy.v1.bidon");
        }
    }
}

pub mod com {
    pub mod iabtechlab {
        pub mod openrtb {
            pub mod v3 {
                tonic::include_proto!("com.iabtechlab.openrtb.v3");
            }
        }
        pub mod adcom {
            pub mod v1 {
                tonic::include_proto!("com.iabtechlab.adcom.v1.context");
                tonic::include_proto!("com.iabtechlab.adcom.v1.enums");
                tonic::include_proto!("com.iabtechlab.adcom.v1.media");
                tonic::include_proto!("com.iabtechlab.adcom.v1.placement");
            }
        }
    }
}


#[cfg(feature = "server")]
pub mod server;

#[cfg(feature = "server")]
pub mod context;

pub mod models;

#[cfg(feature = "server")]
pub(crate) mod header;

#[cfg(feature = "server")]
mod auction;

#[cfg(feature = "server")]
mod controllers {
    pub mod auction;
    pub mod adapter;
}

#[cfg(feature = "server")]
mod bidon_version;
