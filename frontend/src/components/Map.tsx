import { FC, useCallback, useEffect, useState } from "react";
import MapboxGL, { Camera, MapView } from "@rnmapbox/maps";
import { StyleProp, View, ViewStyle } from "react-native";
import { debounce } from "../utils/debounce";
import { getBoundingBox } from "../utils/map";

MapboxGL.setAccessToken(process.env.EXPO_PUBLIC_MAPBOX_ACCESS_TOKEN!);

export const Map: FC<{
  radius: number;
  centerCoordinate: [number, number];
  children?: JSX.Element;
  aspectRatio?: number;
  style?: StyleProp<ViewStyle>;
}> = ({ radius, centerCoordinate, aspectRatio = 1.8, style, children }) => {
  const [bounds, setBounds] = useState(
    getBoundingBox(centerCoordinate, radius, aspectRatio)
  );

  const debouncedSetBounds = useCallback(
    debounce((radius: number) => {
      const bounds = getBoundingBox(centerCoordinate, radius, aspectRatio);
      setBounds(bounds);
    }, 40),
    [centerCoordinate]
  );

  useEffect(() => {
    debouncedSetBounds(radius);
  }, [radius, debouncedSetBounds]);

  return (
    <View
      style={[
        {
          width: "100%",
          aspectRatio,
          borderRadius: 10,
          overflow: "hidden",
        },
        style,
      ]}
    >
      <MapView
        style={{ flex: 1 }}
        logoEnabled={false}
        scaleBarEnabled={false}
        zoomEnabled={false}
        scrollEnabled={false}
        attributionEnabled={false}
      >
        <Camera bounds={bounds} animationDuration={100} />
        {children}
      </MapView>
    </View>
  );
};
