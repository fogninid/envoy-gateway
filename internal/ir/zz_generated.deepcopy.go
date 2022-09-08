//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package ir

import (
	"github.com/envoyproxy/gateway/api/config/v1alpha1"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DirectResponse) DeepCopyInto(out *DirectResponse) {
	*out = *in
	if in.Body != nil {
		in, out := &in.Body, &out.Body
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DirectResponse.
func (in *DirectResponse) DeepCopy() *DirectResponse {
	if in == nil {
		return nil
	}
	out := new(DirectResponse)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HTTPListener) DeepCopyInto(out *HTTPListener) {
	*out = *in
	if in.Hostnames != nil {
		in, out := &in.Hostnames, &out.Hostnames
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = new(TLSListenerConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.Routes != nil {
		in, out := &in.Routes, &out.Routes
		*out = make([]*HTTPRoute, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(HTTPRoute)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HTTPListener.
func (in *HTTPListener) DeepCopy() *HTTPListener {
	if in == nil {
		return nil
	}
	out := new(HTTPListener)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HTTPPathModifier) DeepCopyInto(out *HTTPPathModifier) {
	*out = *in
	if in.FullReplace != nil {
		in, out := &in.FullReplace, &out.FullReplace
		*out = new(string)
		**out = **in
	}
	if in.PrefixMatchReplace != nil {
		in, out := &in.PrefixMatchReplace, &out.PrefixMatchReplace
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HTTPPathModifier.
func (in *HTTPPathModifier) DeepCopy() *HTTPPathModifier {
	if in == nil {
		return nil
	}
	out := new(HTTPPathModifier)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HTTPRoute) DeepCopyInto(out *HTTPRoute) {
	*out = *in
	if in.PathMatch != nil {
		in, out := &in.PathMatch, &out.PathMatch
		*out = new(StringMatch)
		(*in).DeepCopyInto(*out)
	}
	if in.HeaderMatches != nil {
		in, out := &in.HeaderMatches, &out.HeaderMatches
		*out = make([]*StringMatch, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(StringMatch)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.QueryParamMatches != nil {
		in, out := &in.QueryParamMatches, &out.QueryParamMatches
		*out = make([]*StringMatch, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(StringMatch)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.DirectResponse != nil {
		in, out := &in.DirectResponse, &out.DirectResponse
		*out = new(DirectResponse)
		(*in).DeepCopyInto(*out)
	}
	if in.Redirect != nil {
		in, out := &in.Redirect, &out.Redirect
		*out = new(Redirect)
		(*in).DeepCopyInto(*out)
	}
	if in.Destinations != nil {
		in, out := &in.Destinations, &out.Destinations
		*out = make([]*RouteDestination, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(RouteDestination)
				**out = **in
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HTTPRoute.
func (in *HTTPRoute) DeepCopy() *HTTPRoute {
	if in == nil {
		return nil
	}
	out := new(HTTPRoute)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Infra) DeepCopyInto(out *Infra) {
	*out = *in
	if in.Proxy != nil {
		in, out := &in.Proxy, &out.Proxy
		*out = new(ProxyInfra)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Infra.
func (in *Infra) DeepCopy() *Infra {
	if in == nil {
		return nil
	}
	out := new(Infra)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InfraMetadata) DeepCopyInto(out *InfraMetadata) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InfraMetadata.
func (in *InfraMetadata) DeepCopy() *InfraMetadata {
	if in == nil {
		return nil
	}
	out := new(InfraMetadata)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ListenerPort) DeepCopyInto(out *ListenerPort) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ListenerPort.
func (in *ListenerPort) DeepCopy() *ListenerPort {
	if in == nil {
		return nil
	}
	out := new(ListenerPort)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProxyInfra) DeepCopyInto(out *ProxyInfra) {
	*out = *in
	if in.Metadata != nil {
		in, out := &in.Metadata, &out.Metadata
		*out = new(InfraMetadata)
		(*in).DeepCopyInto(*out)
	}
	if in.Config != nil {
		in, out := &in.Config, &out.Config
		*out = new(v1alpha1.EnvoyProxy)
		(*in).DeepCopyInto(*out)
	}
	if in.Listeners != nil {
		in, out := &in.Listeners, &out.Listeners
		*out = make([]ProxyListener, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProxyInfra.
func (in *ProxyInfra) DeepCopy() *ProxyInfra {
	if in == nil {
		return nil
	}
	out := new(ProxyInfra)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProxyListener) DeepCopyInto(out *ProxyListener) {
	*out = *in
	if in.Ports != nil {
		in, out := &in.Ports, &out.Ports
		*out = make([]ListenerPort, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProxyListener.
func (in *ProxyListener) DeepCopy() *ProxyListener {
	if in == nil {
		return nil
	}
	out := new(ProxyListener)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Redirect) DeepCopyInto(out *Redirect) {
	*out = *in
	if in.Scheme != nil {
		in, out := &in.Scheme, &out.Scheme
		*out = new(string)
		**out = **in
	}
	if in.Hostname != nil {
		in, out := &in.Hostname, &out.Hostname
		*out = new(string)
		**out = **in
	}
	if in.Path != nil {
		in, out := &in.Path, &out.Path
		*out = new(HTTPPathModifier)
		(*in).DeepCopyInto(*out)
	}
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(uint32)
		**out = **in
	}
	if in.StatusCode != nil {
		in, out := &in.StatusCode, &out.StatusCode
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Redirect.
func (in *Redirect) DeepCopy() *Redirect {
	if in == nil {
		return nil
	}
	out := new(Redirect)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StringMatch) DeepCopyInto(out *StringMatch) {
	*out = *in
	if in.Exact != nil {
		in, out := &in.Exact, &out.Exact
		*out = new(string)
		**out = **in
	}
	if in.Prefix != nil {
		in, out := &in.Prefix, &out.Prefix
		*out = new(string)
		**out = **in
	}
	if in.SafeRegex != nil {
		in, out := &in.SafeRegex, &out.SafeRegex
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StringMatch.
func (in *StringMatch) DeepCopy() *StringMatch {
	if in == nil {
		return nil
	}
	out := new(StringMatch)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TLSListenerConfig) DeepCopyInto(out *TLSListenerConfig) {
	*out = *in
	if in.ServerCertificate != nil {
		in, out := &in.ServerCertificate, &out.ServerCertificate
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	if in.PrivateKey != nil {
		in, out := &in.PrivateKey, &out.PrivateKey
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TLSListenerConfig.
func (in *TLSListenerConfig) DeepCopy() *TLSListenerConfig {
	if in == nil {
		return nil
	}
	out := new(TLSListenerConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Xds) DeepCopyInto(out *Xds) {
	*out = *in
	if in.HTTP != nil {
		in, out := &in.HTTP, &out.HTTP
		*out = make([]*HTTPListener, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(HTTPListener)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Xds.
func (in *Xds) DeepCopy() *Xds {
	if in == nil {
		return nil
	}
	out := new(Xds)
	in.DeepCopyInto(out)
	return out
}
