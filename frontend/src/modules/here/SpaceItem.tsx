import "react-native-gesture-handler";
import { View, TouchableOpacity } from "react-native";
import { BottomTabNavigationProp } from "@react-navigation/bottom-tabs";
import { FC } from "react";
import { Space, TabsParamList } from "../../types";
import { Text } from "../../components/Text";

export const SpaceItem: FC<{
  data: Space;
  navigation: BottomTabNavigationProp<TabsParamList, "Here", undefined>;
}> = ({ data, navigation }) => {
  return (
    <View
      style={[
        {
          width: "33.33333%",
          padding: 5,
          aspectRatio: 1,
        },
      ]}
    >
      <TouchableOpacity
        style={{
          backgroundColor: `#${data.themeColorHexaCode}`,
          flex: 1,
          borderRadius: 7,
          marginVertical: 0,
          paddingHorizontal: 0,
        }}
        onPress={() => {
          navigation.navigate("Space" as any, { spaceId: data.id });
        }}
      >
        <Text>{JSON.stringify(data)}</Text>
      </TouchableOpacity>
    </View>
  );
};
