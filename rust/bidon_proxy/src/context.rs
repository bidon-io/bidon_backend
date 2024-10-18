use futures::future::BoxFuture;
use hyper::header::HeaderName;
use hyper::{Error, Request, Response, StatusCode, service::Service};
use url::form_urlencoded;
use std::default::Default;
use std::io;
use std::marker::PhantomData;
use std::task::{Poll, Context};
use swagger::auth::{AuthData, Authorization, Bearer, Scopes};
use swagger::{EmptyContext, Has, Pop, Push, XSpanIdString};
use crate::{Api, AuthenticationApi};
use log::error;
use crate::bidon_version::XBidonVersionString;


swagger::new_context_type!(MyContext, MyEmpContext, Option<AuthData>, Option<Authorization>, XSpanIdString, XBidonVersionString);

// Define the BidonContext type
pub(crate) type BidonContext = swagger::make_context_ty!(MyContext, MyEmpContext,  XSpanIdString, XBidonVersionString, Option<AuthData>, Option<Authorization>);

pub struct MakeAddContext<T, A> {
    inner: T,
    marker: PhantomData<A>,
}

impl<T, A, B, C, D,E> MakeAddContext<T, A>
where
    A: Default + Push<XSpanIdString, Result = B>,
    B: Push<XBidonVersionString, Result = C>,
    C: Push<Option<AuthData>, Result = D>,
    D: Push<Option<Authorization>, Result = E>,

{
    pub fn new(inner: T) -> MakeAddContext<T, A> {
        MakeAddContext {
            inner,
            marker: PhantomData,
        }
    }
}

// Make a service that adds context.
impl<Target, T, A, B, C, D,E> Service<Target> for
    MakeAddContext<T, A>
where
    Target: Send,
    A: Default + Push<XSpanIdString, Result = B> + Send,
    B: Push<XBidonVersionString, Result = C>,
    C: Push<Option<AuthData>, Result = D>,
    D: Push<Option<Authorization>, Result = E>,
    E: Send + 'static,
    T: Service<Target> + Send,
    T::Future: Send + 'static
{
    type Error = T::Error;
    type Response = AddContext<T::Response, A, B, C, D,E>;
    type Future = BoxFuture<'static, Result<Self::Response, Self::Error>>;

    fn poll_ready(&mut self, cx: &mut Context<'_>) -> Poll<Result<(), Self::Error>> {
        self.inner.poll_ready(cx)
    }

    fn call(&mut self, target: Target) -> Self::Future {
        let service = self.inner.call(target);

        Box::pin(async move {
            Ok(AddContext::new(service.await?))
        })
    }
}

/// Middleware to add context data from the request
pub struct AddContext<T, A, B, C, D, E>
where
    A: Default + Push<XSpanIdString, Result = B>,
    B: Push<XBidonVersionString, Result = C>,
    C: Push<Option<AuthData>, Result = D>,
    D: Push<Option<Authorization>, Result = E>
{
    inner: T,
    marker: PhantomData<A>,
}

impl<T, A, B, C, D,E> AddContext<T, A, B, C, D, E>
where
    A: Default + Push<XSpanIdString, Result = B>,
    B: Push<XBidonVersionString, Result = C>,
    C: Push<Option<AuthData>, Result = D>,
    D: Push<Option<Authorization>, Result = E>,
{
    pub fn new(inner: T) -> Self {
        AddContext {
            inner,
            marker: PhantomData,
        }
    }
}

impl<T, A, B, C, D, E, ReqBody> Service<Request<ReqBody>> for AddContext<T, A, B, C, D, E>
    where
        A: Default + Push<XSpanIdString, Result=B>,
        B: Push<XBidonVersionString, Result=C>,
        C: Push<Option<AuthData>, Result=D>,
        D: Push<Option<Authorization>, Result=E>,
        E: Send + 'static,
        T: Service<(Request<ReqBody>, E)> + AuthenticationApi
{
    type Response = T::Response;
    type Error = T::Error;
    type Future = T::Future;

    fn poll_ready(&mut self, cx: &mut Context<'_>) -> Poll<Result<(), Self::Error>> {
        self.inner.poll_ready(cx)
    }


    fn call(&mut self, request: Request<ReqBody>) -> Self::Future {
        let context = A::default().push(XSpanIdString::get_or_generate(&request));
        let context = context.push(XBidonVersionString::get_or_generate(&request));
        let headers = request.headers();


        let context = context.push(None::<AuthData>);
        let context = context.push(None::<Authorization>);

        self.inner.call((request, context))
    }
}
