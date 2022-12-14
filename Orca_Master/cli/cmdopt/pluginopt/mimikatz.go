package pluginopt

var MimikatzOptions = []string{
	"privilege::debug",
	"privilege::driver",
	"privilege::security",
	"privilege::tcb",
	"privilege::backup",
	"privilege::restore",
	"privilege::sysenv",
	"privilege::id",
	"privilege::name",
	"sekurlsa::msv",
	"sekurlsa::wdigest",
	"sekurlsa::tspkg",
	"sekurlsa::livessp",
	"sekurlsa::ssp",
	"sekurlsa::logonpasswords",
	"sekurlsa::process",
	"sekurlsa::minidump",
	"sekurlsa::bootkey",
	"sekurlsa::pth",
	"sekurlsa::krbtgt",
	"sekurlsa::dpapisystem",
	"sekurlsa::trust",
	"sekurlsa::backupkeys",
	"sekurlsa::tickets",
	"sekurlsa::ekeys",
	"sekurlsa::dpapi",
	"sekurlsa::credman",
	"lsadump::sam",
	"lsadump::secrets",
	"lsadump::cache",
	"lsadump::lsa",
	"lsadump::trust",
	"lsadump::backupkeys",
	"lsadump::rpdata",
	"lsadump::dcsync",
	"lsadump::dcshadow",
	"lsadump::setntlm",
	"lsadump::changentlm",
	"lsadump::netsync",
	"lsadump::packages",
	"lsadump::mbc",
	"token::whoami",
	"token::list",
	"token::elevate",
	"token::run",
	"token::revert",
	"kerberos::ptt",
	"kerberos::list",
	"kerberos::ask",
	"kerberos::tgt",
	"kerberos::purge",
	"kerberos::golden",
	"kerberos::hash",
	"kerberos::ptc",
	"kerberos::clist",
	"misc::cmd",
	"misc::regedit",
	"misc::taskmgr",
	"misc::ncroutemon",
	"misc::detours",
	"misc::memssp",
	"misc::skeleton",
	"misc::compressme",
	"misc::lock",
	"misc::wp",
	"misc::mflt",
	"misc::easyntlmchall",
	"misc::clip",
	"crypto::providers",
	"crypto::stores",
	"crypto::certificates",
	"crypto::keys",
	"crypto::sc",
	"crypto::hash",
	"crypto::system",
	"crypto::scauth",
	"crypto::certtohw",
	"crypto::capi",
	"crypto::cng",
	"crypto::extract",
	"crypto::kutil",
	"process::list",
	"process::exports",
	"process::imports",
	"process::start",
	"process::stop",
	"process::suspend",
	"process::resume",
	"process::run",
	"process::runp",
	"ts::multirdp",
	"ts::sessions",
	"ts::remote",
	"event::drop",
	"event::clear",
}
