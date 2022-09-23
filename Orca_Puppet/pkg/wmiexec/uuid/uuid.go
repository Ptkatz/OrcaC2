package uuid

type UUID [16]byte

var NULL = UUID{}

var NDRTransferSyntax_V2 = fromStringInternalOnly("8a885d04-1ceb-11c9-9fe8-08002b104860")

//https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-rpce/dca648a5-42d3-432c-9927-2f22e50fa266
var NDR64TransferSyntax = fromStringInternalOnly("71710533-beba-4937-8319-b5dbef9ccc36")

var BindTimeFeatureReneg = fromStringInternalOnly("6cb71c2c-9812-4540-0300-000000000000")

//https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-dcom/c25391af-f59e-40da-885e-cc84076673e4
var IID_IRemUnknown2 = fromStringInternalOnly("00000143-0000-0000-c000-000000000046")
var IID_IActivationPropertiesIn = fromStringInternalOnly("000001a2-0000-0000-c000-000000000046")
var CLSID_ActivationPropertiesIn = fromStringInternalOnly("00000338-0000-0000-c000-000000000046")
var CLSID_SpecialSystemProperties = fromStringInternalOnly("000001b9-0000-0000-c000-000000000046")
var CLSID_InstantiationInfo = fromStringInternalOnly("000001ab-0000-0000-c000-000000000046")
var CLSID_ActivationContextInfo = fromStringInternalOnly("000001a5-0000-0000-c000-000000000046")
var CLSID_SecurityInfo = fromStringInternalOnly("000001a6-0000-0000-c000-000000000046")
var CLSID_ServerLocationInfo = fromStringInternalOnly("000001a4-0000-0000-c000-000000000046")
var CLSID_ScmRequestInfo = fromStringInternalOnly("000001aa-0000-0000-c000-000000000046")
var IID_IObjectExporter = fromStringInternalOnly("99fcfec4-5260-101b-bbcb-00aa0021347a")
var IID_IContext = fromStringInternalOnly("000001c0-0000-0000-c000-000000000046")
var CLSID_ContextMarshaler = fromStringInternalOnly("0000033b-0000-0000-c000-000000000046")
var IID_IRemoteSCMActivator = fromStringInternalOnly("000001A0-0000-0000-C000-000000000046")

//https://answers.microsoft.com/en-us/windows/forum/all/the-server-8bc3f05e-d86b-11d0-a075-00c04fb68820/7500c1d2-b873-4e68-af8c-89fe7e848658
var CLSID_WMIAppID = fromStringInternalOnly("8bc3f05e-d86b-11d0-a075-00c04fb68820")

//https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-wmi/3485541f-6950-4e6d-98cb-1ed4bb143441
//https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-wmi/485026a6-d7e0-4ef8-a44f-43e5853fff9d
var CLSID_WbemLevel1Login = fromStringInternalOnly("f309ad18-d86a-11d0-a075-00c04fb68820")
var IID_IWbemLoginClientID = fromStringInternalOnly("d4781cd6-e5d3-44df-ad94-930efe48a887")
var IID_IWbemServices = fromStringInternalOnly("9556dc99-828c-11cf-a37e-00aa003240c7")
