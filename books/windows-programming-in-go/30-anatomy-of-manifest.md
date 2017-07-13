---
Id: wpig-30
Title: Anatomy of .manifest file
Format: Markdown
Tags: go
CreatedAt: 2017-07-12
PublishedOn: 2017-09-12
Collection: go-windows
Description: What is .manifest file and how to use it.
Status: invisible
---

A manifest is an [XML file](https://msdn.microsoft.com/en-us/library/windows/desktop/aa375365(v=vs.85).aspx) that tells Windows some important information about your program.

There are 2 ways to provide manifest file:
* as a separate file alongside your executable. If your program is `foo.exe`, manifest file should be named `foo.exe.manifest`
* embedded as a resource inside `.exe`files

Here's a `.manifest` file that I use most often:
```xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
  <assemblyIdentity version="1.0.0.0" processorArchitecture="*" name="SomeFunkyNameHere" type="win32"/>
  <dependency>
    <dependentAssembly>
      <assemblyIdentity type="win32" name="Microsoft.Windows.Common-Controls" version="6.0.0.0" processorArchitecture="*" publicKeyToken="6595b64144ccf1df" language="*"/>
    </dependentAssembly>
  </dependency>
  <application xmlns="urn:schemas-microsoft-com:asm.v3">
    <windowsSettings>
      <dpiAware xmlns="http://schemas.microsoft.com/SMI/2005/WindowsSettings">True</dpiAware>
    </windowsSettings>
  </application>
  <compatibility xmlns="urn:schemas-microsoft-com:compatibility.v1">
    <application>
      <!-- Windows Vista -->
      <supportedOS Id="{e2011457-1546-43c5-a5fe-008deee3d3f0}"/>
      <!-- Windows 7 -->
      <supportedOS Id="{35138b9a-5d96-4fbd-8e2d-a2440225f93a}"/>
      <!-- Windows 8 -->
      <supportedOS Id="{4a2f28e3-53b9-4441-ba9c-d69d4a4a6e38}"/>
      <!-- Windows 8.1 -->
      <supportedOS Id="{1f676c76-80e1-4239-95bb-83d0f6d0da78}"/>
      <!-- Windows 10 -->
      <supportedOS Id="{8e0f7a12-bfb3-4fe8-b9a5-48fd50a15a9a}"/>
    </application>
  </compatibility>
</assembly>
```

Let's deconstruct it's content.

## Opting into newest version of common controls

Standard UI elements like text boxes, list boxes etc. are called common controls.

The OS library that implements common controls keeps improving, adding new functionality and potentially subtly changing the behavior.

Programs written for older version of Windows, e.g. XP, might not be tested with latest version of common controls code (e.g. on Windows 10) and therefore have subtle bugs.

Microsoft really cares about backwards compatibility so it ships old versions of common controls code.

Windows assumes that unless you explicitly state that you support latest version of common contols then you don't and Windows will have your program use older version.

Latest common controls look better so you should tell Windows you want them. You do it by including the following in manifest:

```xml
<dependentAssembly>
    <assemblyIdentity type="win32" name="Microsoft.Windows.Common-Controls" version="6.0.0.0" processorArchitecture="*" publicKeyToken="6595b64144ccf1df" language="*"/>
</dependentAssembly>
```

## Confirming you support the latest versions of the OS

Software that works just fine on XP might work incorrectly on Windows 10 due to doing things that were never officially documented but just happened to work on XP.

Microsoft cares about compatibility so they have "compatibility hacks" code which e.g. emulates XP behavior on Windows 10 in order to allow buggy programs to work. If they detect that program is mis-behaving they'll enable compatibility mode for that program.

You should test your code on Windows 10 and explicitly tell Windows that you did that and therefore don't want compatibility hacks for your program. You do that by including this in manifest file:

```xml
<compatibility xmlns="urn:schemas-microsoft-com:compatibility.v1">
  <application>
    <!-- Windows Vista -->
    <supportedOS Id="{e2011457-1546-43c5-a5fe-008deee3d3f0}"/>
    <!-- Windows 7 -->
    <supportedOS Id="{35138b9a-5d96-4fbd-8e2d-a2440225f93a}"/>
    <!-- Windows 8 -->
    <supportedOS Id="{4a2f28e3-53b9-4441-ba9c-d69d4a4a6e38}"/>
    <!-- Windows 8.1 -->
    <supportedOS Id="{1f676c76-80e1-4239-95bb-83d0f6d0da78}"/>
    <!-- Windows 10 -->
    <supportedOS Id="{8e0f7a12-bfb3-4fe8-b9a5-48fd50a15a9a}"/>
  </application>
</compatibility>
```

## Opting out of DPI scaling

For a long time most displays had very similar resolution density, close to 96 dpi (dots, or pixels, per inch).

That meant that if you opened a window 800x600 pixels in size, it's physical size on the monitor would be the same on most displays.

That changed in the last few years and we got much higher resolution displays. I write it on Microsoft Surface Laptop which has DPI of XXX, which is XX times higher than 96 dpi.

What it means is that now 800x600 window looks very tiny.

This is not a great behavior so by default newer versions of Windows will automatically scale the window size based on display's dpi, as relative to 96 dpi.

Unfortunately that results in blurry display because scaling is imperfect.

To avoid bluriness you have to tell Windows you don't want automatic DPI scaling. To do that, add the following to the manifest:

```xml
  <application xmlns="urn:schemas-microsoft-com:asm.v3">
    <windowsSettings>
      <dpiAware xmlns="http://schemas.microsoft.com/SMI/2005/WindowsSettings">True</dpiAware>
    </windowsSettings>
  </application>
```

This enables system DPI scaling. There's also `True/PM` version which enables per-monitor awareness, which is more complicated to support.

Opting out of DPI scaling means that you'll have to scale all the sizes yourself. We'll cover that in **TODO: link to DPI scaling**

## Running as admin

Some operations (like writing to protected directories like `c:\Program Files`) are not allowed for regular users. One has to be an admin to do them. On personal laptops the user is admin but for security is not running with administrative priviledges.

If you want your program to run with administrative priviledges by default, you need to add this to the manifest:

```xml
<trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
  <security>
    <requestedPrivileges>
      <requestedExecutionLevel level="requireAdministrator" uiAccess="false"/>
    </requestedPrivileges>
  </security>
</trustInfo>
```

This is useful for programs like installers, which might want to write to `c:\Program Files`.

When user runs such program, Windows will ask the user if he allows running with admin priviledges. This is known as UAC (User Access Control) prompt.

Another technique is to start running as non-admin and elevate to admin mode by re-launching the program.

That would delay UAC prompt to when you need admin priviledges (as opposed to showing it when program starts).

