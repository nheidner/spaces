import "dotenv/config";

export default () => ({
  expo: {
    name: "spaces-p",
    slug: "spaces-p",
    scheme: "spaces-p",
    version: "1.0.0",
    orientation: "portrait",
    icon: "src/assets/icon.png",
    userInterfaceStyle: "dark",
    splash: {
      image: "src/assets/splash.png",
      resizeMode: "contain",
    },
    plugins: [
      "expo-apple-authentication",
      [
        "expo-location",
        {
          locationWhenInUsePermission:
            "Allow $(PRODUCT_NAME) to use your location.",
        },
      ],
      [
        "@rnmapbox/maps",
        {
          RNMapboxMapsDownloadToken:
            process.env.EXPO_PUBLIC_MAPBOX_ACCESS_TOKEN,
        },
      ],
    ],
    assetBundlePatterns: ["**/*"],
    ios: {
      supportsTablet: true,
      bundleIdentifier: "com.anonymous.spacesp",
      usesAppleSignIn: true,
    },
    android: {
      adaptiveIcon: {
        foregroundImage: "src/assets/adaptive-icon.png",
        backgroundColor: "#ffffff",
      },
      package: "com.anonymous.spacesp",
      permissions: [
        "android.permission.ACCESS_COARSE_LOCATION",
        "android.permission.ACCESS_FINE_LOCATION",
        "android.permission.FOREGROUND_SERVICE",
      ],
    },
    web: {
      favicon: "src/assets/favicon.png",
    },
    extra: {
      eas: {
        projectId: "e689bcf2-50ed-4a2d-a994-a49ec5fde1c8",
      },
    },
  },
});
