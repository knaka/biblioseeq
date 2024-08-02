// Language Guide (proto 3) | Protocol Buffers Documentation https://protobuf.dev/programming-guides/proto3/

// @generated by protoc-gen-connect-es v0.11.0 with parameter "target=js+dts"
// @generated from file v1/main.proto (package v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { ContentRequest, ContentResponse, CurrentTimeRequest, CurrentTimeResponse, OpenFileRequest, OpenFileResponse, OpenURLRequest, OpenURLResponse, QueryRequest, QueryResponse, StatusRequest, StatusResponse, VersionInfoRequest, VersionInfoResponse } from "./main_pb.js";
import { MethodKind } from "@bufbuild/protobuf";

/**
 * @generated from service v1.MainService
 */
export const MainService = {
  typeName: "v1.MainService",
  methods: {
    /**
     * @generated from rpc v1.MainService.VersionInfo
     */
    versionInfo: {
      name: "VersionInfo",
      I: VersionInfoRequest,
      O: VersionInfoResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc v1.MainService.Status
     */
    status: {
      name: "Status",
      I: StatusRequest,
      O: StatusResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc v1.MainService.CurrentTime
     */
    currentTime: {
      name: "CurrentTime",
      I: CurrentTimeRequest,
      O: CurrentTimeResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc v1.MainService.Query
     */
    query: {
      name: "Query",
      I: QueryRequest,
      O: QueryResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc v1.MainService.Content
     */
    content: {
      name: "Content",
      I: ContentRequest,
      O: ContentResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc v1.MainService.OpenFile
     */
    openFile: {
      name: "OpenFile",
      I: OpenFileRequest,
      O: OpenFileResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc v1.MainService.OpenURL
     */
    openURL: {
      name: "OpenURL",
      I: OpenURLRequest,
      O: OpenURLResponse,
      kind: MethodKind.Unary,
    },
  }
};

