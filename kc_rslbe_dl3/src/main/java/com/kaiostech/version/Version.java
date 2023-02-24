package com.kaiostech.version;

import java.lang.StringBuilder;

public class Version {
    public static String VERSION = "1.0.0";
    public static String BUILDID = "49c14dadbe1d40d0ded949a71ee68329071c38c0-20230224001718";
    public static String DATE = "Fri Feb 24 00:17:18 PST 2023";
    public static String BUILDER = "maen";
    public static String HOSTNAME = "localhost.localdomain";
    public static String KERNEL_VERSION = "Linux localhost.localdomain 3.10.0-1160.81.1.el7.x86_64 #1 SMP Fri Dec 16 17:29:43 UTC 2022 x86_64 x86_64 x86_64 GNU/Linux";
    public static String KERNEL_RELEASE = "CentOS Linux release 7.9.2009 (Core)";

    public static String Get_Version() {
        return VERSION;
    }

    public static String Get_BuildID()  {
       return BUILDID;
    }

    public static String Get_Build_Date() {
       return DATE;
    }

    public static String Get_Builder() {
       return BUILDER;
    }

    public static String Get_Hostname() {
       return HOSTNAME;
    }

    public static String getFullVersion() {
       StringBuilder buffer;
       buffer = new StringBuilder();
 
       buffer.append("        Version: ").append(VERSION).append("\n");
       buffer.append("       Build ID: ").append(BUILDID).append("\n");
       buffer.append("     Build Date: ").append(DATE).append("\n");
       buffer.append("        Builder: ").append(BUILDER).append("\n");
       buffer.append("     Build Host: ").append(HOSTNAME).append("\n");
       buffer.append(" Kernel Version: ").append(KERNEL_VERSION).append("\n");
       buffer.append(" Kernel Release: ").append(KERNEL_RELEASE).append("\n");
       buffer.append("\n");

       return buffer.toString();
    }
    
}


