use serde_json::{json, Value};

pub fn get_auction_request() -> Value {
    json!( {
        "app": {
        "key": "string",
        "bundle": "string",
        "framework": "string",
        "version": "string",
        "framework_version": "string",
        "plugin_version": "string",
        "skadn": [
        "string"
        ],
        "sdk_version": "string"
    },
        "device": {
        "geo": {
            "lat": 0,
            "lon": 0,
            "accuracy": 0,
            "lastfix": 0,
            "country": "string",
            "city": "string",
            "zip": "string",
            "utcoffset": 0
        },
        "ua": "string",
        "make": "string",
        "model": "string",
        "os": "string",
        "osv": "string",
        "hwv": "string",
        "h": 0,
        "w": 0,
        "ppi": 0,
        "pxratio": 0,
        "js": 0,
        "language": "string",
        "carrier": "string",
        "mccmnc": "string",
        "connection_type": "ETHERNET",
        "type": "PHONE"
    },
        "ext": "string",
        "geo": {
        "lat": 0,
        "lon": 0,
        "accuracy": 0,
        "lastfix": 0,
        "country": "string",
        "city": "string",
        "zip": "string",
        "utcoffset": 0
    },
        "regs": {
        "coppa": true,
        "gdpr": true,
        "us_privacy": "string",
        "eu_privacy": "string",
        "iab": {}
    },
        "session": {
        "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
        "launch_ts": 0,
        "launch_monotonic_ts": 0,
        "start_ts": 0,
        "start_monotonic_ts": 0,
        "ts": 0,
        "monotonic_ts": 0,
        "memory_warnings_ts": [
        0
        ],
        "memory_warnings_monotonic_ts": [
        0
        ],
        "ram_used": 0,
        "ram_size": 0,
        "storage_free": 0,
        "storage_used": 0,
        "battery": 0,
        "cpu_usage": 0
    },
        "segment": {
        "id": "string",
        "uid": "string",
        "ext": "string"
    },
        "token": "string",
        "user": {
        "idfa": "c013dc70-1e9a-4f92-96cd-e867b31c4b7d",
        "tracking_authorization_status": "string",
        "idfv": "2b12f8b8-5c4f-441e-9d03-b7cef54f5689",
        "idg": "ca403851-fd71-4229-9ee2-418852281604",
        "consent": {},
        "coppa": true
    },
        "adapters": {
        "property1": {
            "version": "string",
            "sdk_version": "string"
        },
        "property2": {
            "version": "string",
            "sdk_version": "string"
        }
    },
        "ad_object": {
        "auction_id": "string",
        "auction_key": "string",
        "auction_configuration_id": 0,
        "auction_configuration_uid": "string",
        "auction_pricefloor": 0,
        "orientation": "PORTRAIT",
        "demands": {
            "meta": {
                "token": "token",
                "status": "SUCCESS",
                "token_finish_ts": 1700000000000u64,
                "token_start_ts": 1700000000001u64
            },
            "vungle": {
                "token": "token",
                "status": "FAIL",
                "token_finish_ts": 1700000000000u64,
                "token_start_ts": 1700000000001u64
            }
        },
        "banner": {
            "format": "BANNER"
        },
        "interstitial": {},
        "rewarded": {}
    },
        "test": true,
        "tmax": 0
    })
}
