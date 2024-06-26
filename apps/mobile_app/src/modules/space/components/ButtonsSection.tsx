import React, { FC, useState } from "react";
import { ScrollView, Share, View } from "react-native";
import { template } from "../../../styles/template";
import { Text } from "../../../components/Text";
import { SpaceStackParamList, Uuid } from "../../../types";
import { PressableTransformation } from "../../../components/PressableTransformation";
import { StackNavigationProp } from "@react-navigation/stack";
import { useNavigation } from "@react-navigation/native";
import { hexToRgb } from "../../../utils/hex_to_rgb";
import { useNotification } from "../../../utils/notifications";
import { useSpaceColor } from "../../../hooks/use_space_color";

export const ButtonsSection: FC = () => {
  const navigation = useNavigation<StackNavigationProp<SpaceStackParamList>>();
  const [joined, setJoined] = useState(false);
  const notification = useNotification();

  const handleJoin = () => {
    if (joined) {
      notification.showNotification({
        title: "Left",
        description: "You have left the space",
        type: "info",
      });
    } else {
      notification.showNotification({
        title: "Joined",
        description: "You have joined the space",
        type: "success",
      });
    }

    setJoined((oldJoined) => !oldJoined);
  };

  const handleShare = async () => {
    try {
      await Share.share({
        message: "Check out this space",
      });
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <ScrollView
      contentContainerStyle={{
        flexDirection: "row",
        gap: 10,
        paddingBottom: 10,
      }}
      style={{ backgroundColor: template.colors.white }}
      horizontal
      showsHorizontalScrollIndicator={false}
    >
      {[
        {
          text: joined ? "subscribed" : "subscribe",
          onPress: handleJoin,
          isActive: joined,
        },
        {
          text: "share",
          onPress: handleShare,
        },
        {
          text: "info",
          onPress: () => {
            navigation.navigate("Info");
          },
        },
      ].map((data, index) => {
        return (
          <Button
            key={index}
            text={data.text}
            onPress={data.onPress}
            isActive={data.isActive}
          />
        );
      })}
    </ScrollView>
  );
};

const Button: FC<{ text: string; onPress: () => void; isActive?: boolean }> = ({
  text,
  onPress,
  isActive = false,
}) => {
  const spaceColor = useSpaceColor();

  return (
    <PressableTransformation onPress={onPress} hitSlop={10}>
      <View
        style={{
          paddingHorizontal: 20,
          paddingVertical: 7,
          backgroundColor: isActive
            ? template.colors.grayLightBackground
            : hexToRgb(spaceColor, 0.15),
          borderRadius: 999,
        }}
      >
        <Text
          style={{
            fontWeight: template.fontWeight.bold,
            fontSize: 16,
            color: isActive ? template.colors.textLight : spaceColor,
          }}
        >
          {text}
        </Text>
      </View>
    </PressableTransformation>
  );
};
