import { template } from "../styles/template";
import { NotificationType } from "../types";
import { FC, useEffect } from "react";
import { ActivityIndicator, Text, View, StyleSheet } from "react-native";
import { useNotificationState } from "../components/context/NotificationContext";
import Animated, {
  useAnimatedStyle,
  useSharedValue,
  withTiming,
} from "react-native-reanimated";
import { useSafeAreaInsets } from "react-native-safe-area-context";

const backgroundColors: Record<NotificationType, string> = {
  info: template.colors.purple,
  error: template.colors.error,
  success: template.colors.success,
  loading: template.colors.purpleLight,
};

// TODO: make safe area insets work

export const CustomNotification: FC<{
  title: string;
  description: string;
  type: NotificationType;
}> = () => {
  const [notificationState] = useNotificationState();
  const heightSv = useSharedValue(0);

  const insets = useSafeAreaInsets();

  const typeSv = useSharedValue(notificationState?.type);

  useEffect(() => {
    typeSv.value = notificationState?.type;
  }, [notificationState?.type]);

  useEffect(() => {
    const numberElements = [
      notificationState?.title,
      notificationState?.description,
    ].filter((a) => !!a).length;

    heightSv.value = numberElements * 22 + 24;
  }, [notificationState?.title, notificationState?.description]);

  const animatedBackgroundColor = useAnimatedStyle(() => {
    return {
      backgroundColor: withTiming(backgroundColors[typeSv.value || "info"], {
        duration: 100,
      }),
    };
  });

  const animatedHeight = useAnimatedStyle(() => {
    return {
      height: withTiming(heightSv.value, { duration: 200 }),
    };
  });

  const loadingIndicator =
    notificationState?.type === "loading" ? (
      <ActivityIndicator
        color={template.colors.white}
        style={{ marginRight: 10 }}
      />
    ) : null;

  const title = notificationState?.title ? (
    <Text style={styles.title}>{notificationState.title}</Text>
  ) : null;

  const description = notificationState?.description ? (
    <Text style={styles.description}>{notificationState.description}</Text>
  ) : null;

  return (
    <Animated.View
      style={[
        styles.container,
        { paddingTop: insets.top },
        animatedBackgroundColor,
      ]}
    >
      <Animated.View style={[styles.innerContainer, animatedHeight]}>
        {loadingIndicator}
        <View>
          {title}
          {description}
        </View>
      </Animated.View>
    </Animated.View>
  );
};

const styles = StyleSheet.create({
  container: { flex: 1 },
  innerContainer: {
    width: "100%",
    flexDirection: "row",
    justifyContent: "center",
    alignItems: "center",
    height: "auto",
  },
  title: {
    color: template.colors.white,
    textAlign: "center",
    fontWeight: template.fontWeight.bold,
    fontSize: 16,
  },
  description: {
    color: template.colors.white,
    textAlign: "center",
    marginTop: 3,
  },
});
