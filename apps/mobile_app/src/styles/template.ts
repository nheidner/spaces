import { TextStyle } from "react-native";

export const template = {
  colors: {
    gray: "#d5d5d5",
    purpleLightBackground: "#C081FF33",
    purpleLight: "#C081FF",
    purple: "#A347FF",
    grayLightBackground: "#F5F5F5",
    text: "#444",
    textLight: "#888",
    white: "#fff",
    error: "#FD7979",
    red: "#AB0303",
    success: "#4CAF50",
  },
  fontWeight: {
    bold: "700" as TextStyle["fontWeight"],
  },
  paddings: {
    md: 15,
    screen: 25,
  },
  margins: {
    sm: 15,
    md: 25,
  },
  borderRadius: {
    md: 8,
    screen: 7,
  },
  height: {
    header: 56,
  },
  fontSizes: {
    sm: 12,
    md: 16,
    lg: 20,
  },
};
