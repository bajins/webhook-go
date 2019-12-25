1>1/* ::
:: by bajins https://www.bajins.com

@echo off
md "%~dp0$testAdmin$" 2>nul
if not exist "%~dp0$testAdmin$" (
    echo bajins���߱�����Ŀ¼��д��Ȩ��! >&2
    exit /b 1
) else rd "%~dp0$testAdmin$"

:: �����ӳٻ���������չ
setlocal enabledelayedexpansion

:: %~f0 ��ʾ��ǰ������ľ���·��,ȥ�����ŵ�����·��
cscript -nologo -e:jscript "%~f0" get7z
if not "%errorlevel%" == "0" (
    @cmd /k
    goto :EXIT
)

:: ��Ҫ������ļ����ļ��и�Ŀ¼
set root=%~dp0
:: ��Ҫ������ļ����ļ���
set files=data

:: ���� %0 ���䵽һ��·��
set currentPath=%~p0
:: �滻\Ϊ,�ţ�Ҳ�����滻Ϊ�ո�
set currentPath=%currentPath:\=,%
:: ˳��ѭ�����������һ��Ϊ��ǰĿ¼
for %%a in (%currentPath%) do set CurrentDirectoryName=%%a
:: �����ɵ��ļ�����ǰһ����
set project=%CurrentDirectoryName%
:: �����ɵ��ļ�������һ���֣���ǰһ���ֽ������
set allList=_darwin_386,_darwin_amd64,_freebsd_386,_freebsd_amd64,_freebsd_arm,
set allList=%allList%_netbsd_386,_netbsd_amd64,_netbsd_arm,_openbsd_386,
set allList=%allList%_openbsd_amd64,_windows_386.exe,_windows_amd64.exe,
set allList=%allList%_linux_386,_linux_amd64,_linux_arm,_linux_mips,
set allList=%allList%_linux_mips64,_linux_mips64le,_linux_mipsle,_linux_s390x

:GETGOX
set GOPROXY=https://goproxy.io
go get github.com/mitchellh/gox

for %%i in (%allList%) do (
    :: ����������ļ������������´��
    if not exist "%project%%%i" (
        gox
        if not %errorlevel% == 0 (
            goto :GETGOX
        )
        :: ɾ���ɵ�ѹ�����ļ�
        del *.zip *.tar *.gz
    )
)


:: ʹ��7zѹ��
for %%i in (%allList%) do (
    set runFile=%project%%%i
    :: !!Ϊsetlocal EnableDelayedExpansionȡ������ֵ
    if exist "!runFile!" (
        :: �жϱ����ַ������Ƿ�����ַ���
        echo %%i | findstr linux >nul && (
            :: ��7zѹ����tar
            7za a -ttar %project%%%i.tar %files% !runFile!
            :: ��7z��tarѹ����gz
            7za a -tgzip %project%%%i.tar.gz %project%%%i.tar
            :: ɾ��tar�ļ��Ͷ������ļ�
            del *.tar !runFile!
            
        ) || (
            :: ��7zѹ���ļ�Ϊzip
            7za a %project%%%i.zip %files% !runFile!
            :: ɾ���������ļ�
            del !runFile!
        )
    )
)



goto :EXIT

:EXIT
:: �����ӳٻ���������չ������ִ��
endlocal&exit /b %errorlevel%
*/

// ****************************  JavaScript  *******************************


var Argv = WScript.Arguments;
for (i = 0; i < Argv.length; i++) {
    WScript.StdOut.WriteLine("������" + Argv(i));
}

if (Argv.length > 0) {
    switch (Argv(0)) {
        case "get7z":
            try{
                get7z();
            }catch(e){
                WScript.StdErr.WriteLine(e.message);
                // �������˳�
                WScript.Quit(1);
            }
            break;
        default:
            help();
    }
    // �����˳�
    WScript.Quit(0);
}


/**
 * HTTP����
 *
 * @param method        GET,POST
 * @param url           �����ַ
 * @param dataType      "",text,stream,xml,json
 * @param data          ���ݣ�{key:value}��ʽ
 * @param contentType   ���͵��������ͣ�multipart/form-data��
 * application/x-www-form-urlencoded��Ĭ�ϣ���text/plain
 * @returns {string|Document|any}
 */
function request(method, url, dataType, data, contentType) {
    if (url == "" || url == null || url.length <= 0) {
        throw new Error("����url����Ϊ�գ�");
    }
    if (method == "" || method == null || method.length <= 0) {
        method = "GET";
    } else {
        // ���ַ���ת��Ϊ��д
        method = method.toUpperCase();
    }
    if (contentType == "" || contentType == null || contentType.length <= 0) {
        contentType = "application/x-www-form-unlenconded;charset=utf-8";
    }
    var XMLHTTPVersions = [
        'WinHttp.WinHttpRequest.5.1',
        'WinHttp.WinHttpRequest.5.0',
        'Msxml2.ServerXMLHTTP.6.0',
        'Msxml2.ServerXMLHTTP.5.0',
        'Msxml2.ServerXMLHTTP.4.0',
        'Msxml2.ServerXMLHTTP.3.0',
        'Msxml2.ServerXMLHTTP',
        'MSXML2.XMLHTTP.6.0',
        'MSXML2.XMLHTTP.5.0',
        'MSXML2.XMLHTTP.4.0',
        'MSXML2.XMLHTTP.3.0',
        'MSXML2.XMLHTTP',
        'Microsoft.XMLHTTP'
    ];
    var XMLHTTP;
    for (var i = 0; i < XMLHTTPVersions.length; i++) {
        try {
            XMLHTTP = new ActiveXObject(XMLHTTPVersions[i]);
            break;
        } catch (e) {
            WScript.StdOut.Write(XMLHTTPVersions[i]);
            WScript.StdOut.WriteLine("��" + e.message);
        }
    }

    //������ת����Ϊquerystring��ʽ
    var paramarray = [];
    for (key in data) {
        paramarray.push(key + "=" + data[key]);
    }
    var params = paramarray.join("&");

    switch (method) {
        case "POST":
            // 0�첽��1ͬ��
            XMLHTTP.Open(method, url, 0);
            XMLHTTP.SetRequestHeader("CONTENT-TYPE", contentType);
            XMLHTTP.Send(params);
            break;
        default:
            // Ĭ��GET����
            if (params == "" || params.length == 0 || params == null) {
                // 0�첽��1ͬ��
                XMLHTTP.Open(method, url, 0);
            } else {
                XMLHTTP.Open(method, url + "?" + params, 0);
            }
            XMLHTTP.SetRequestHeader("CONTENT-TYPE", contentType);
            XMLHTTP.Send();
    }

    // ���ַ���ת��ΪСд
    dataType = dataType.toLowerCase();
    switch (dataType) {
        case "text":
            return XMLHTTP.responseText;
            break;
        case "stream":
            return XMLHTTP.responseStream;
            break;
        case "xml":
            return XMLHTTP.responseXML;
            break;
        case "json":
            return eval("(" + XMLHTTP.responseText + ")");
            break;
        default:
            return XMLHTTP.responseBody;
    }
}


/**
 * �����ļ�
 *
 * @param url
 * @param directory �ļ��洢Ŀ¼
 * @param filename  �ļ�����Ϊ��Ĭ�Ͻ�ȡurl�е��ļ���
 * @returns {string}
 */
function download(url, directory, filename) {
    if (url == "" || url == null || url.length <= 0) {
        throw new Error("����url����Ϊ�գ�");
    }
    if (directory == "" || directory == null || directory.length <= 0) {
        throw new Error("�ļ��洢Ŀ¼����Ϊ�գ�");
    }

    var fso = new ActiveXObject("Scripting.FileSystemObject");
    // ���Ŀ¼������
    if (!fso.FolderExists(directory)) {
        // ����Ŀ¼
        var strFolderName = fso.CreateFolder(directory);
    }

    if (filename == "" || filename == null || filename.length <= 0) {
        filename = url.substring(url.lastIndexOf("/") + 1);
        // ȥ���ļ�����������ţ�����֮ǰ�ģ��ַ�
        filename = filename.replace(/^.*(\&|\=|\?|\/)/ig, "");
    }
    var path = directory + "\\" + filename;

    var ADO = new ActiveXObject("ADODB.Stream");
    ADO.Mode = 3;
    ADO.Type = 1;
    ADO.Open();
    ADO.Write(request("GET", url, ""));
    ADO.SaveToFile(path, 2);
    ADO.Close();

    // ����ļ�������
    if (!fso.FileExists(path)) {
        return "";
    }
    return path;
}

/**
 * ��ȡ7-Zip
 *
 */
function get7z() {
    var shell = new ActiveXObject("WScript.shell");
    // ִ��7z�����ж��Ƿ�ִ�гɹ�
    var out = shell.Run("cmd /c 7za", 0, true);
    var directory = "c:\\windows";
    var url = "https://github.com/woytu/woytu.github.io/releases/download/v1.0/7za.exe";
    // ���ִ��ʧ��˵��7z������
    if (out == 1) {
        download(url, directory);
    }
    // ִ��7z�����ж��Ƿ�ִ�гɹ�
    out = shell.Run("cmd /c 7za", 0, true);
    var fso = new ActiveXObject("Scripting.FileSystemObject");
    // ���ִ��ʧ�ܣ������ļ�������
    if (out == 1 || !fso.FileExists(directory + "\\7za.exe")) {
        get7z();
    }
}