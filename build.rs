//! This file is used to compile the proto files in the proto directory

use prost_build::Config;

fn main() {
    let proto_dirs = ["proto/adcom/proto", "proto/openrtb/proto", "proto/proto"];

    // Print the directories for debugging (optional)
    println!("Protobuf directories: {:?}", proto_dirs);

    // Compile the proto files, including all found directories for imports
    tonic_build::configure()
        .codec_path("crate::codec::ProstRegistryCodec")
        .compile_protos(&["org/bidon/proto/v1/services.proto"], &proto_dirs)
        .unwrap_or_else(|e| panic!("Failed to compile protos: {:?}", e));

    Config::new()
        .compile_protos(
            &[
                "org/bidon/proto/v1/mediation/mediation.proto",
                "org/bidon/proto/v1/context/context.proto",
                "com/iabtechlab/openrtb/v3/openrtb.proto",
                "com/iabtechlab/adcom/v1/adcom.proto",
                "com/iabtechlab/adcom/v1/context/context.proto",
                "com/iabtechlab/adcom/v1/enums/enums.proto",
                "com/iabtechlab/adcom/v1/media/media.proto",
                "com/iabtechlab/adcom/v1/placement/placement.proto",
            ],
            &proto_dirs,
        )
        .unwrap_or_else(|e| panic!("Failed to compile protos: {:?}", e));
}
