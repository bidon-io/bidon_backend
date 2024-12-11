use std::{
    marker::PhantomData,
    sync::{atomic::AtomicPtr, OnceLock},
};

use prost::{ExtensionRegistry, Message};
use tonic::{
    codec::{BufferSettings, Codec, DecodeBuf, Decoder, EncodeBuf, Encoder},
    Status,
};

/// A [`Codec`] that implements `application/grpc+proto` via the prost
/// library with a extension registry.
#[derive(Debug, Clone, Default, Copy)]
pub struct ProstRegistryCodec<T, U> {
    _pd: PhantomData<(T, U)>,
}

static REGISTRY: OnceLock<AtomicPtr<ExtensionRegistry>> = OnceLock::new();

/// Initialize the global extension registry.
/// This must be called before any decoding operations.
///
/// # Safety
/// This function is unsafe because it stores a raw pointer to the registry.
/// The caller must ensure that the registry outlives all uses of the codec.
pub unsafe fn init_registry(registry: ExtensionRegistry) {
    let ptr = Box::into_raw(Box::new(registry));
    REGISTRY.get_or_init(|| AtomicPtr::new(ptr));
}

impl<T, U> ProstRegistryCodec<T, U>
where
    T: Message + Send + 'static,
    U: Message + Default + Send + 'static,
{
    /// A tool for building custom codecs based on prost encoding and decoding.
    /// See the codec_buffers example for one possible way to use this.
    pub fn raw_encoder(buffer_settings: BufferSettings) -> <Self as Codec>::Encoder {
        ProstEncoder {
            _pd: PhantomData,
            buffer_settings,
        }
    }

    /// A tool for building custom codecs based on prost encoding and decoding.
    /// See the codec_buffers example for one possible way to use this.
    pub fn raw_decoder(buffer_settings: BufferSettings) -> <Self as Codec>::Decoder {
        ProstDecoder {
            _pd: PhantomData,
            buffer_settings,
        }
    }
}

impl<T, U> Codec for ProstRegistryCodec<T, U>
where
    T: Message + Send + 'static,
    U: Message + Default + Send + 'static,
{
    type Encode = T;
    type Decode = U;

    type Encoder = ProstEncoder<T>;
    type Decoder = ProstDecoder<U>;

    fn encoder(&mut self) -> Self::Encoder {
        ProstEncoder {
            _pd: PhantomData,
            buffer_settings: BufferSettings::default(),
        }
    }

    fn decoder(&mut self) -> Self::Decoder {
        ProstDecoder {
            _pd: PhantomData,
            buffer_settings: BufferSettings::default(),
        }
    }
}

/// A [`Encoder`] that knows how to encode `T`.
#[derive(Debug, Clone, Default)]
pub struct ProstEncoder<T> {
    _pd: PhantomData<T>,
    buffer_settings: BufferSettings,
}

impl<T> ProstEncoder<T> {
    /// Get a new encoder with explicit buffer settings
    pub fn new(buffer_settings: BufferSettings) -> Self {
        Self {
            _pd: PhantomData,
            buffer_settings,
        }
    }
}

impl<T: Message> Encoder for ProstEncoder<T> {
    type Item = T;
    type Error = Status;

    fn encode(&mut self, item: Self::Item, buf: &mut EncodeBuf<'_>) -> Result<(), Self::Error> {
        item.encode(buf)
            .expect("Message only errors if not enough space");

        Ok(())
    }

    fn buffer_settings(&self) -> BufferSettings {
        self.buffer_settings
    }
}

/// A [`Decoder`] that knows how to decode `U`.
#[derive(Debug, Clone, Default)]
pub struct ProstDecoder<U> {
    _pd: PhantomData<U>,
    buffer_settings: BufferSettings,
}

impl<U> ProstDecoder<U> {
    /// Get a new decoder with explicit buffer settings
    pub fn new(buffer_settings: BufferSettings) -> Self {
        Self {
            _pd: PhantomData,
            buffer_settings,
        }
    }
}

impl<U: Message + Default> Decoder for ProstDecoder<U> {
    type Item = U;
    type Error = Status;
    fn decode(&mut self, buf: &mut DecodeBuf<'_>) -> Result<Option<Self::Item>, Self::Error> {
        let item = if let Some(registry_ptr) = REGISTRY.get() {
            Message::decode_with_extensions(buf, unsafe {
                &*registry_ptr.load(std::sync::atomic::Ordering::Acquire)
            })
        } else {
            Message::decode(buf)
        }
        .map(Option::Some)
        .map_err(from_decode_error)?;

        Ok(item)
    }

    fn buffer_settings(&self) -> BufferSettings {
        self.buffer_settings
    }
}

fn from_decode_error(error: prost::DecodeError) -> Status {
    // Map Protobuf parse errors to an INTERNAL status code, as per
    // https://github.com/grpc/grpc/blob/master/doc/statuscodes.md
    Status::internal(error.to_string())
}
