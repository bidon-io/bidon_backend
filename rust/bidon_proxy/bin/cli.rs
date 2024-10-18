//! CLI tool driving the API client
use anyhow::{anyhow, Context, Result};
use log::{debug, info};
// models may be unused if all inputs are primitive types
#[allow(unused_imports)]
use bidon::{
    models, ApiNoContext, Client, ContextWrapperExt,
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
use simple_logger::SimpleLogger;
use structopt::StructOpt;
use swagger::{AuthData, ContextBuilder, EmptyContext, Push, XSpanIdString};

type ClientContext = swagger::make_context_ty!(
    ContextBuilder,
    EmptyContext,
    Option<AuthData>,
    XSpanIdString
);

#[derive(StructOpt, Debug)]
#[structopt(
    name = "SDK API",
    version = "1.0.0",
    about = "CLI access to SDK API"
)]
struct Cli {
    #[structopt(subcommand)]
    operation: Operation,

    /// Address or hostname of the server hosting this API, including optional port
    #[structopt(short = "a", long, default_value = "http://localhost")]
    server_address: String,

    /// Path to the client private key if using client-side TLS authentication
    #[cfg(not(any(target_os = "macos", target_os = "windows", target_os = "ios")))]
    #[structopt(long, requires_all(&["client-certificate", "server-certificate"]))]
    client_key: Option<String>,

    /// Path to the client's public certificate associated with the private key
    #[cfg(not(any(target_os = "macos", target_os = "windows", target_os = "ios")))]
    #[structopt(long, requires_all(&["client-key", "server-certificate"]))]
    client_certificate: Option<String>,

    /// Path to CA certificate used to authenticate the server
    #[cfg(not(any(target_os = "macos", target_os = "windows", target_os = "ios")))]
    #[structopt(long)]
    server_certificate: Option<String>,

    /// If set, write output to file instead of stdout
    #[structopt(short, long)]
    output_file: Option<String>,

    #[structopt(flatten)]
    verbosity: clap_verbosity_flag::Verbosity,
}

#[derive(StructOpt, Debug)]
enum Operation {
    /// Get config
    GetConfig {
        /// Version of the Bidon SDK
        x_bidon_version: String,
        /// Config request
        #[structopt(parse(try_from_str = parse_json))]
        config_request: models::ConfigRequest,
    },
    /// Get OpenAPI specification
    GetOpenApiSpec {
    },
    /// Auction
    GetAuction {
        /// Version of the Bidon SDK
        x_bidon_version: String,
        /// Ad type
        #[structopt(parse(try_from_str = parse_json))]
        ad_type: models::GetAuctionAdTypeParameter,
        /// Auction request
        #[structopt(parse(try_from_str = parse_json))]
        auction_request: models::AuctionRequest,
    },
    /// Click
    PostClick {
        /// Version of the Bidon SDK
        x_bidon_version: String,
        /// Ad type
        #[structopt(parse(try_from_str = parse_json))]
        ad_type: models::GetAuctionAdTypeParameter,
        /// Click request
        #[structopt(parse(try_from_str = parse_json))]
        click_request: models::ClickRequest,
    },
    /// Loss
    PostLoss {
        /// Version of the Bidon SDK
        x_bidon_version: String,
        /// Ad type
        #[structopt(parse(try_from_str = parse_json))]
        ad_type: models::GetAuctionAdTypeParameter,
        /// Loss request
        #[structopt(parse(try_from_str = parse_json))]
        loss_request: models::LossRequest,
    },
    /// Reward
    PostReward {
        /// Version of the Bidon SDK
        x_bidon_version: String,
        /// Ad type for the reward request
        #[structopt(parse(try_from_str = parse_json))]
        ad_type: models::PostRewardAdTypeParameter,
        /// Reward request
        #[structopt(parse(try_from_str = parse_json))]
        reward_request: models::RewardRequest,
    },
    /// Show
    PostShow {
        /// Version of the Bidon SDK
        x_bidon_version: String,
        /// Ad type
        #[structopt(parse(try_from_str = parse_json))]
        ad_type: models::GetAuctionAdTypeParameter,
        /// Show request
        #[structopt(parse(try_from_str = parse_json))]
        show_request: models::ShowRequest,
    },
    /// Stats
    PostStats {
        /// Version of the Bidon SDK
        x_bidon_version: String,
        /// Ad type
        #[structopt(parse(try_from_str = parse_json))]
        ad_type: models::GetAuctionAdTypeParameter,
        /// Stats request
        #[structopt(parse(try_from_str = parse_json))]
        stats_request: models::StatsRequest,
    },
    /// Win
    PostWin {
        /// Version of the Bidon SDK
        x_bidon_version: String,
        /// Ad type
        #[structopt(parse(try_from_str = parse_json))]
        ad_type: models::GetAuctionAdTypeParameter,
        /// Win request
        #[structopt(parse(try_from_str = parse_json))]
        win_request: models::WinRequest,
    },
}

#[cfg(not(any(target_os = "macos", target_os = "windows", target_os = "ios")))]
fn create_client(args: &Cli, context: ClientContext) -> Result<Box<dyn ApiNoContext<ClientContext>>> {
    if args.client_certificate.is_some() {
        debug!("Using mutual TLS");
        let client = Client::try_new_https_mutual(
            &args.server_address,
            args.server_certificate.clone().unwrap(),
            args.client_key.clone().unwrap(),
            args.client_certificate.clone().unwrap(),
        )
        .context("Failed to create HTTPS client")?;
        Ok(Box::new(client.with_context(context)))
    } else if args.server_certificate.is_some() {
        debug!("Using TLS with pinned server certificate");
        let client =
            Client::try_new_https_pinned(&args.server_address, args.server_certificate.clone().unwrap())
                .context("Failed to create HTTPS client")?;
        Ok(Box::new(client.with_context(context)))
    } else {
        debug!("Using client without certificates");
        let client =
            Client::try_new(&args.server_address).context("Failed to create HTTP(S) client")?;
        Ok(Box::new(client.with_context(context)))
    }
}

#[cfg(any(target_os = "macos", target_os = "windows", target_os = "ios"))]
fn create_client(args: &Cli, context: ClientContext) -> Result<Box<dyn ApiNoContext<ClientContext>>> {
    let client =
        Client::try_new(&args.server_address).context("Failed to create HTTP(S) client")?;
    Ok(Box::new(client.with_context(context)))
}

#[tokio::main]
async fn main() -> Result<()> {
    let args = Cli::from_args();
    if let Some(log_level) = args.verbosity.log_level() {
        SimpleLogger::new().with_level(log_level.to_level_filter()).init()?;
    }

    debug!("Arguments: {:?}", &args);

    let auth_data: Option<AuthData> = None;

    #[allow(trivial_casts)]
    let context = swagger::make_context!(
        ContextBuilder,
        EmptyContext,
        auth_data,
        XSpanIdString::default()
    );

    let client = create_client(&args, context)?;

    let result = match args.operation {
        Operation::GetConfig {
            x_bidon_version,
            config_request,
        } => {
            info!("Performing a GetConfig request");

            let result = client.get_config(
                x_bidon_version,
                config_request,
            ).await?;
            debug!("Result: {:?}", result);

            match result {
                GetConfigResponse::ConfigResponse
                (body)
                => "ConfigResponse\n".to_string()
                   +
                    &serde_json::to_string_pretty(&body)?,
            }
        }
        Operation::GetOpenApiSpec {
        } => {
            info!("Performing a GetOpenApiSpec request");

            let result = client.get_open_api_spec(
            ).await?;
            debug!("Result: {:?}", result);

            match result {
                GetOpenApiSpecResponse::OpenAPIJSONSpecification
                (body)
                => "OpenAPIJSONSpecification\n".to_string()
                   +
                    &serde_json::to_string_pretty(&body)?,
                GetOpenApiSpecResponse::Error
                (body)
                => "Error\n".to_string()
                   +
                    &serde_json::to_string_pretty(&body)?,
            }
        }
        Operation::GetAuction {
            x_bidon_version,
            ad_type,
            auction_request,
        } => {
            info!("Performing a GetAuction request on {:?}", (
                &ad_type
            ));

            let result = client.get_auction(
                x_bidon_version,
                ad_type,
                auction_request,
            ).await?;
            debug!("Result: {:?}", result);

            match result {
                GetAuctionResponse::AuctionResponse
                (body)
                => "AuctionResponse\n".to_string()
                   +
                    &serde_json::to_string_pretty(&body)?,
            }
        }
        Operation::PostClick {
            x_bidon_version,
            ad_type,
            click_request,
        } => {
            info!("Performing a PostClick request on {:?}", (
                &ad_type
            ));

            let result = client.post_click(
                x_bidon_version,
                ad_type,
                click_request,
            ).await?;
            debug!("Result: {:?}", result);

            match result {
                PostClickResponse::ClickResponse
                (body)
                => "ClickResponse\n".to_string()
                   +
                    &serde_json::to_string_pretty(&body)?,
            }
        }
        Operation::PostLoss {
            x_bidon_version,
            ad_type,
            loss_request,
        } => {
            info!("Performing a PostLoss request on {:?}", (
                &ad_type
            ));

            let result = client.post_loss(
                x_bidon_version,
                ad_type,
                loss_request,
            ).await?;
            debug!("Result: {:?}", result);

            match result {
                PostLossResponse::LossResponse
                (body)
                => "LossResponse\n".to_string()
                   +
                    &serde_json::to_string_pretty(&body)?,
            }
        }
        Operation::PostReward {
            x_bidon_version,
            ad_type,
            reward_request,
        } => {
            info!("Performing a PostReward request on {:?}", (
                &ad_type
            ));

            let result = client.post_reward(
                x_bidon_version,
                ad_type,
                reward_request,
            ).await?;
            debug!("Result: {:?}", result);

            match result {
                PostRewardResponse::RewardResponse
                (body)
                => "RewardResponse\n".to_string()
                   +
                    &serde_json::to_string_pretty(&body)?,
            }
        }
        Operation::PostShow {
            x_bidon_version,
            ad_type,
            show_request,
        } => {
            info!("Performing a PostShow request on {:?}", (
                &ad_type
            ));

            let result = client.post_show(
                x_bidon_version,
                ad_type,
                show_request,
            ).await?;
            debug!("Result: {:?}", result);

            match result {
                PostShowResponse::ShowResponse
                (body)
                => "ShowResponse\n".to_string()
                   +
                    &serde_json::to_string_pretty(&body)?,
            }
        }
        Operation::PostStats {
            x_bidon_version,
            ad_type,
            stats_request,
        } => {
            info!("Performing a PostStats request on {:?}", (
                &ad_type
            ));

            let result = client.post_stats(
                x_bidon_version,
                ad_type,
                stats_request,
            ).await?;
            debug!("Result: {:?}", result);

            match result {
                PostStatsResponse::StatsResponse
                (body)
                => "StatsResponse\n".to_string()
                   +
                    &serde_json::to_string_pretty(&body)?,
            }
        }
        Operation::PostWin {
            x_bidon_version,
            ad_type,
            win_request,
        } => {
            info!("Performing a PostWin request on {:?}", (
                &ad_type
            ));

            let result = client.post_win(
                x_bidon_version,
                ad_type,
                win_request,
            ).await?;
            debug!("Result: {:?}", result);

            match result {
                PostWinResponse::WinResponse
                (body)
                => "WinResponse\n".to_string()
                   +
                    &serde_json::to_string_pretty(&body)?,
            }
        }
    };

    if let Some(output_file) = args.output_file {
        std::fs::write(output_file, result)?
    } else {
        println!("{}", result);
    }
    Ok(())
}

// May be unused if all inputs are primitive types
#[allow(dead_code)]
fn parse_json<'a, T: serde::de::Deserialize<'a>>(json_string: &'a str) -> Result<T> {
    serde_json::from_str(json_string).map_err(|err| anyhow!("Error parsing input: {}", err))
}
