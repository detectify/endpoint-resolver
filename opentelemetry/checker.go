// Code generated by gowrap. DO NOT EDIT.
// template: https://raw.githubusercontent.com/hexdigest/gowrap/6c8f05695fec23df85903a8da0af66ac414e2a63/templates/opentelemetry
// gowrap: http://github.com/hexdigest/gowrap

package opentelemetry

//go:generate gowrap gen -p github.com/detectify/endpoint-resolver -i Checker -t https://raw.githubusercontent.com/hexdigest/gowrap/6c8f05695fec23df85903a8da0af66ac414e2a63/templates/opentelemetry -o checker.go -l ""

import (
	"context"

	endpointresolver "github.com/detectify/endpoint-resolver"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// CheckerWithTracing implements endpointresolver.Checker interface instrumented with opentracing spans
type CheckerWithTracing struct {
	endpointresolver.Checker
	_instance      string
	_spanDecorator func(span trace.Span, params, results map[string]interface{})
}

// NewCheckerWithTracing returns CheckerWithTracing
func NewCheckerWithTracing(base endpointresolver.Checker, instance string, spanDecorator ...func(span trace.Span, params, results map[string]interface{})) CheckerWithTracing {
	d := CheckerWithTracing{
		Checker:   base,
		_instance: instance,
	}

	if len(spanDecorator) > 0 && spanDecorator[0] != nil {
		d._spanDecorator = spanDecorator[0]
	}

	return d
}

// CheckAllPorts implements endpointresolver.Checker
func (_d CheckerWithTracing) Ports(ctx context.Context, ips []string, ports []int) (openPorts []int, err error) {
	ctx, _span := otel.Tracer(_d._instance).Start(ctx, "endpointresolver.Checker.Ports")
	defer func() {
		if _d._spanDecorator != nil {
			_d._spanDecorator(_span, map[string]interface{}{
				"ctx":   ctx,
				"ips":   ips,
				"ports": ports}, map[string]interface{}{
				"openPorts": openPorts,
				"err":       err})
		} else if err != nil {
			_span.RecordError(err)
			_span.SetAttributes(
				attribute.String("event", "error"),
				attribute.String("message", err.Error()),
			)
		}

		_span.End()
	}()
	return _d.Checker.Ports(ctx, ips, ports)
}

// DNSCheckWithExternalProvider implements endpointresolver.Checker
func (_d CheckerWithTracing) ExternalDNS(ctx context.Context, hostname string, externalDNS []string) (err error) {
	ctx, _span := otel.Tracer(_d._instance).Start(ctx, "endpointresolver.Checker.ExternalDNS")
	defer func() {
		if _d._spanDecorator != nil {
			_d._spanDecorator(_span, map[string]interface{}{
				"ctx":         ctx,
				"hostname":    hostname,
				"externalDNS": externalDNS}, map[string]interface{}{
				"err": err})
		} else if err != nil {
			_span.RecordError(err)
			_span.SetAttributes(
				attribute.String("event", "error"),
				attribute.String("message", err.Error()),
			)
		}

		_span.End()
	}()
	return _d.Checker.ExternalDNS(ctx, hostname, externalDNS)
}

// HTTPCheck implements endpointresolver.Checker
func (_d CheckerWithTracing) HTTP(ctx context.Context, userAgent string, hostname string, customHeaders map[string]string, openPorts []int) (sa1 []string, err error) {
	ctx, _span := otel.Tracer(_d._instance).Start(ctx, "endpointresolver.Checker.HTTP")
	defer func() {
		if _d._spanDecorator != nil {
			_d._spanDecorator(_span, map[string]interface{}{
				"ctx":           ctx,
				"userAgent":     userAgent,
				"hostname":      hostname,
				"customHeaders": customHeaders,
				"openPorts":     openPorts}, map[string]interface{}{
				"sa1": sa1,
				"err": err})
		} else if err != nil {
			_span.RecordError(err)
			_span.SetAttributes(
				attribute.String("event", "error"),
				attribute.String("message", err.Error()),
			)
		}

		_span.End()
	}()
	return _d.Checker.HTTP(ctx, userAgent, hostname, customHeaders, openPorts)
}

// NativeDNSCheck implements endpointresolver.Checker
func (_d CheckerWithTracing) NativeDNS(ctx context.Context, hostname string) (ips []string, err error) {
	ctx, _span := otel.Tracer(_d._instance).Start(ctx, "endpointresolver.Checker.NativeDNS")
	defer func() {
		if _d._spanDecorator != nil {
			_d._spanDecorator(_span, map[string]interface{}{
				"ctx":      ctx,
				"hostname": hostname}, map[string]interface{}{
				"ips": ips,
				"err": err})
		} else if err != nil {
			_span.RecordError(err)
			_span.SetAttributes(
				attribute.String("event", "error"),
				attribute.String("message", err.Error()),
			)
		}

		_span.End()
	}()
	return _d.Checker.NativeDNS(ctx, hostname)
}
