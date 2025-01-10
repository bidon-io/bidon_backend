use axum::{
    http::StatusCode,
    response::{IntoResponse, Response},
    Json,
};
use derive_more::Display;
use tonic::Status;

#[derive(Debug, Display)]
pub enum BiddingError {
    #[display("gRPC error: {}", _0)]
    GrpcError(Status),
    #[display("HTTP status: {}, error: {}", _0, _1)]
    HttpError(StatusCode, String),
    #[display("Validation error: {}", _0)]
    ValidationError(String),
    #[display("Serialization error: {}", _0)]
    SerializationError(String),
    #[display("Unexpected error: {}", _0)]
    UnexpectedError(String),
}

impl std::error::Error for BiddingError {}

impl BiddingError {
    /// Maps gRPC status codes to appropriate HTTP status codes
    /// Source: https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto
    fn grpc_to_http_status(code: tonic::Code) -> StatusCode {
        match code {
            tonic::Code::Cancelled => StatusCode::from_u16(499).unwrap(),
            tonic::Code::Unknown => StatusCode::INTERNAL_SERVER_ERROR,
            tonic::Code::InvalidArgument => StatusCode::BAD_REQUEST,
            tonic::Code::DeadlineExceeded => StatusCode::GATEWAY_TIMEOUT,
            tonic::Code::NotFound => StatusCode::NOT_FOUND,
            tonic::Code::AlreadyExists => StatusCode::CONFLICT,
            tonic::Code::PermissionDenied => StatusCode::FORBIDDEN,
            tonic::Code::ResourceExhausted => StatusCode::TOO_MANY_REQUESTS,
            tonic::Code::FailedPrecondition => StatusCode::BAD_REQUEST,
            tonic::Code::Aborted => StatusCode::CONFLICT,
            tonic::Code::OutOfRange => StatusCode::BAD_REQUEST,
            tonic::Code::Unimplemented => StatusCode::NOT_IMPLEMENTED,
            tonic::Code::Internal => StatusCode::INTERNAL_SERVER_ERROR,
            tonic::Code::Unavailable => StatusCode::SERVICE_UNAVAILABLE,
            tonic::Code::DataLoss => StatusCode::INTERNAL_SERVER_ERROR,
            tonic::Code::Unauthenticated => StatusCode::UNAUTHORIZED,
            _ => StatusCode::INTERNAL_SERVER_ERROR,
        }
    }
}

impl IntoResponse for BiddingError {
    fn into_response(self) -> Response {
        let (status_code, error_message) = match self {
            BiddingError::GrpcError(status) => {
                match serde_json::from_str::<GrpcErrorMessage>(status.message()) {
                    Ok(e) => {
                        let code = StatusCode::from_u16(e.error.code as u16)
                            .inspect_err(|e| {
                                tracing::debug!("Error parsing gRPC error code: {}", e)
                            })
                            .unwrap_or_else(|_| Self::grpc_to_http_status(status.code()));
                        (code, e.error.message)
                    }
                    Err(e) => {
                        let msg = format!(
                            "Error parsing gRPC error from: '{}', err: {}",
                            status.message(),
                            e
                        );
                        tracing::debug!(msg);

                        (Self::grpc_to_http_status(status.code()), msg)
                    }
                }
            }
            BiddingError::HttpError(status, msg) => {
                (status, format!("Upstream HTTP service error: {}", msg))
            }
            BiddingError::ValidationError(msg) => {
                (StatusCode::BAD_REQUEST, format!("Invalid request: {}", msg))
            }
            BiddingError::SerializationError(msg) => (
                StatusCode::BAD_REQUEST,
                format!("Serialization error: {}", msg),
            ),
            BiddingError::UnexpectedError(msg) => (
                StatusCode::INTERNAL_SERVER_ERROR,
                format!("Internal server error: {}", msg),
            ),
        };

        let error = GrpcErrorMessage {
            error: GrpcErrorDetails {
                code: status_code.as_u16(),
                message: error_message,
            },
        };

        (status_code, Json(error)).into_response()
    }
}

#[derive(serde::Deserialize, serde::Serialize)]
pub struct GrpcErrorMessage {
    #[serde(rename = "error")]
    pub error: GrpcErrorDetails,
}

#[derive(serde::Deserialize, serde::Serialize)]
pub struct GrpcErrorDetails {
    pub code: u16,
    pub message: String,
}

#[cfg(test)]
mod tests {
    use super::*;

    use axum::body::to_bytes;
    use axum::response::Response;
    use serde_json::Value;
    use tokio; // Ensure Tokio is included in your dependencies

    #[tokio::test]
    async fn test_bidding_error_response() {
        // Create the error
        let error = BiddingError::GrpcError(Status::invalid_argument("Invalid input"));

        // Convert the error into an Axum response
        let response: Response = error.into_response();

        // Assert that the status code is BAD_REQUEST (400)
        assert_eq!(response.status(), StatusCode::BAD_REQUEST);

        // Split the response into parts to extract the body
        let (_, body) = response.into_parts();

        // Convert the body into bytes asynchronously
        let bytes = to_bytes(body, usize::MAX)
            .await
            .expect("Failed to read body");

        // Deserialize the bytes into a JSON value
        let body_json: Value = serde_json::from_slice(&bytes).expect("Failed to parse JSON");

        // Assert the JSON structure and values
        assert_eq!(
            body_json["error"]["message"],
            "Error parsing gRPC error from: 'Invalid input', err: expected value at line 1 column 1"
        );
        assert_eq!(body_json["error"]["code"], 400);
    }

    #[tokio::test]
    async fn test_bidding_error_grpc_json_response() {
        // Create a JSON string representing a GrpcError
        let grpc_error_json = r#"{"error": {"message": "Custom error message", "code": 422}}"#;

        // Create the error with the JSON string as the status message
        let error = BiddingError::GrpcError(Status::invalid_argument(grpc_error_json));

        // Convert the error into an Axum response
        let response: Response = error.into_response();

        // Assert that the status code matches the one in our JSON (422)
        assert_eq!(response.status(), StatusCode::UNPROCESSABLE_ENTITY);

        // Split the response into parts to extract the body
        let (_, body) = response.into_parts();

        // Convert the body into bytes asynchronously
        let bytes = to_bytes(body, usize::MAX)
            .await
            .expect("Failed to read body");

        // Deserialize the bytes into a JSON value
        let body_json: Value = serde_json::from_slice(&bytes).expect("Failed to parse JSON");

        // Assert the JSON structure and values
        assert_eq!(body_json["error"]["message"], "Custom error message");
        assert_eq!(body_json["error"]["code"], 422);
    }
}
