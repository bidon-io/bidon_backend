#![allow(missing_docs, unused_variables, trivial_casts)]


#[allow(unused_imports)]
use futures::{future, Stream, stream};
#[allow(unused_imports)]
use bidon::{Api, Claims, Client, models,
            GetConfigResponse,
            GetOpenApiSpecResponse,
            GetAuctionResponse,
            PostClickResponse,
            PostLossResponse,
            PostRewardResponse,
            PostShowResponse,
            PostStatsResponse,
            PostWinResponse,
};
use clap::{App, Arg};

// NOTE: Set environment variable RUST_LOG to the name of the executable (or "cargo run") to activate console logging for all loglevels.
//     See https://docs.rs/env_logger/latest/env_logger/  for more details

#[allow(unused_imports)]
use log::info;

// swagger::Has may be unused if there are no examples
#[allow(unused_imports)]
use swagger::{AuthData, ContextBuilder, EmptyContext, Has, Push, XSpanIdString};
use bidon::bidon_version::XBidonVersionString;
use bidon::context::{MyEmpContext, MyContext};

type ClientContext = swagger::make_context_ty!(MyContext, MyEmpContext, Option<AuthData>, XSpanIdString, XBidonVersionString);

mod client_auth;
use client_auth::build_token;


// rt may be unused if there are no examples
#[allow(unused_mut)]
fn main() {
    env_logger::init();

    let matches = App::new("client")
        .arg(Arg::with_name("operation")
            .help("Sets the operation to run")
            .possible_values(&[
                "GetOpenApiSpec",
            ])
            .required(true)
            .index(1))
        .arg(Arg::with_name("https")
            .long("https")
            .help("Whether to use HTTPS or not"))
        .arg(Arg::with_name("host")
            .long("host")
            .takes_value(true)
            .default_value("localhost")
            .help("Hostname to contact"))
        .arg(Arg::with_name("port")
            .long("port")
            .takes_value(true)
            .default_value("8080")
            .help("Port to contact"))
        .get_matches();

    // Create Bearer-token with a fixed key (secret) for test purposes.
    // In a real (production) system this Bearer token should be obtained via an external Identity/Authentication-server
    // Ensure that you set the correct algorithm and encodingkey that matches what is used on the server side.
    // See https://github.com/Keats/jsonwebtoken for more information
    let auth_token = build_token(
        Claims {
            sub: "tester@acme.com".to_owned(),
            company: "ACME".to_owned(),
            iss: "my_identity_provider".to_owned(),
            // added a very long expiry time
            aud: "org.acme.Resource_Server".to_string(),
            exp: 10000000000,
            // In this example code all available Scopes are added, so the current Bearer Token gets fully authorization.
            scopes:
            "".to_owned(),
        },
        b"secret").unwrap();

    let auth_data = if !auth_token.is_empty() {
        Some(AuthData::Bearer(swagger::auth::Bearer { token: auth_token }))
    } else {
        // No Bearer-token available, so return None
        None
    };

    let is_https = matches.is_present("https");
    let base_url = format!("{}://{}:{}",
                           if is_https { "https" } else { "http" },
                           matches.value_of("host").unwrap(),
                           matches.value_of("port").unwrap());


    let mut client: Box<dyn Api<ClientContext>> = if matches.is_present("https") {
        // Using Simple HTTPS
        let client = Box::new(Client::try_new_https(&base_url)
            .expect("Failed to create HTTPS client"));
        client
    } else {
        // Using HTTP
        let client = Box::new(Client::try_new_http(
            &base_url)
            .expect("Failed to create HTTP client"));
        client
    };

    let mut rt = tokio::runtime::Runtime::new().unwrap();

    match matches.value_of("operation") {
        /* Disabled because there's no example.
        Some("GetConfig") => {
            let result = rt.block_on(client.get_config(
                  "x_bidon_version_example".to_string(),
                  ???
            ));
            info!("{:?} (X-Span-ID: {:?})", result, (client.context() as &dyn Has<XSpanIdString>).get().clone());
        },
        */
        Some("GetOpenApiSpec") => {
            let context: ClientContext =
                swagger::make_context!(MyContext, MyEmpContext, auth_data, XSpanIdString::default(), XBidonVersionString::default());

            let result = rt.block_on(client.get_open_api_spec(&context.clone()));
            info!("{:?} (X-Span-ID: {:?})", result, swagger::Has::<XSpanIdString>::get(&context).clone());
        }
        /* Disabled because there's no example.
        Some("GetAuction") => {
            let result = rt.block_on(client.get_auction(
                  "x_bidon_version_example".to_string(),
                  ???,
                  ???
            ));
            info!("{:?} (X-Span-ID: {:?})", result, (client.context() as &dyn Has<XSpanIdString>).get().clone());
        },
        */
        /* Disabled because there's no example.
        Some("PostClick") => {
            let result = rt.block_on(client.post_click(
                  "x_bidon_version_example".to_string(),
                  ???,
                  ???
            ));
            info!("{:?} (X-Span-ID: {:?})", result, (client.context() as &dyn Has<XSpanIdString>).get().clone());
        },
        */
        /* Disabled because there's no example.
        Some("PostLoss") => {
            let result = rt.block_on(client.post_loss(
                  "x_bidon_version_example".to_string(),
                  ???,
                  ???
            ));
            info!("{:?} (X-Span-ID: {:?})", result, (client.context() as &dyn Has<XSpanIdString>).get().clone());
        },
        */
        /* Disabled because there's no example.
        Some("PostReward") => {
            let result = rt.block_on(client.post_reward(
                  "x_bidon_version_example".to_string(),
                  ???,
                  ???
            ));
            info!("{:?} (X-Span-ID: {:?})", result, (client.context() as &dyn Has<XSpanIdString>).get().clone());
        },
        */
        /* Disabled because there's no example.
        Some("PostShow") => {
            let result = rt.block_on(client.post_show(
                  "x_bidon_version_example".to_string(),
                  ???,
                  ???
            ));
            info!("{:?} (X-Span-ID: {:?})", result, (client.context() as &dyn Has<XSpanIdString>).get().clone());
        },
        */
        /* Disabled because there's no example.
        Some("PostStats") => {
            let result = rt.block_on(client.post_stats(
                  "x_bidon_version_example".to_string(),
                  ???,
                  ???
            ));
            info!("{:?} (X-Span-ID: {:?})", result, (client.context() as &dyn Has<XSpanIdString>).get().clone());
        },
        */
        /* Disabled because there's no example.
        Some("PostWin") => {
            let result = rt.block_on(client.post_win(
                  "x_bidon_version_example".to_string(),
                  ???,
                  ???
            ));
            info!("{:?} (X-Span-ID: {:?})", result, (client.context() as &dyn Has<XSpanIdString>).get().clone());
        },
        */
        _ => {
            panic!("Invalid operation provided")
        }
    }
}
