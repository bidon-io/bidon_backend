//! This file is used to compile the proto files in the proto directory

fn main() {
    let proto_dirs = ["proto/adcom/proto", "proto/openrtb/proto", "proto/proto"];

    // Print the directories for debugging (optional)
    println!("Protobuf directories: {:?}", proto_dirs);

    // Compile the proto files, including all found directories for imports
    tonic_build::configure()
        .compile_protos(
            &["galaxy/v1/services.proto",
                "com/iabtechlab/openrtb/v3/openrtb.proto",
                "com/iabtechlab/adcom/v1/adcom.proto"
            ],  // List of .proto files to compile
            &proto_dirs,                  // Directories where imports may be found
        )
        .unwrap_or_else(|e| panic!("Failed to compile protos: {:?}", e));
}
