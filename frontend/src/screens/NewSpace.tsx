import { KeyboardAvoidingView, View } from "react-native";
import { BottomTabScreenProps } from "@react-navigation/bottom-tabs";
import { FC, useEffect, useState } from "react";
import { TabsParamList } from "../types";
import { useLocation } from "../hooks/use_location";
// import { RnMap } from "../modules/new_space/RnMap";
import { MapboxMap } from "../modules/new_space/MapboxMap";
import { template } from "../styles/template";
import { BottomSheetScrollView } from "@gorhom/bottom-sheet";
import { Header } from "../modules/new_space/Header";
import { Text } from "../components/Text";
import { TextInput } from "../components/form/TextInput";
import { Label } from "../components/form/Label";
import { ColorPicker } from "../modules/new_space/ColorPicker";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import { PrimaryButton } from "../components/form/PrimaryButton";
import { Slider } from "../components/form/Slider";
import { useNewSpaceState } from "../components/context/NewSpaceContext";
import { ZodError, z } from "zod";

const screenPaddingHorizontal = 20;
const gapSize = 10; // This is the uniform gap size you want
const numberOfColumns = 7;

const colors = [
  template.colors.purple,
  "#212078",
  "#69701e",
  template.colors.purple,
  "#ddd",
  "#faa",
  template.colors.purple,
  "#ddd",
  "#faa",
  template.colors.purple,
  "#ddd",
  "#faa",
  template.colors.purple,
  "#ddd",
  "#faa",
  template.colors.purple,
  "#ddd",
];

export const NewSpaceScreen: FC<
  BottomTabScreenProps<TabsParamList, "NewSpace">
> = () => {
  const [spaceNameErrors, setSpaceNameErrors] = useState<string[]>([]);

  const [newSpaceState, dispatch] = useNewSpaceState();
  const { radius, name, selectedColorIndex } = newSpaceState;

  useEffect(() => {
    try {
      z.string().max(40).parse(name);

      setSpaceNameErrors([]);
    } catch (error: ZodError | any) {
      if (error instanceof ZodError) {
        setSpaceNameErrors(error.errors.map((e) => e.message));
      }
    }
  }, [name]);

  const handleRadiusChange = (newRadius: number) => {
    dispatch!({ type: "SET_RADIUS", newRadius });
  };

  const handleNameChange = (newName: string) => {
    dispatch!({ type: "SET_NAME", newName });
  };

  const handleSelectedColorIndexChange = (newIndex: number) => {
    dispatch!({ type: "SELECT_COLOR_INDEX", newIndex });
  };

  const { location, permissionGranted } = useLocation();
  const insets = useSafeAreaInsets();

  const handleSubmit = () => {
    console.log("submitting");
  };

  if (!permissionGranted) {
    return (
      <View>
        <Text>no permission granted</Text>
      </View>
    );
  }

  if (!location) {
    return (
      <View>
        <Text>error</Text>
      </View>
    );
  }

  return (
    <View
      style={{
        flex: 1,
      }}
    >
      <Header />
      <KeyboardAvoidingView
        style={{ flex: 1 }}
        behavior="padding"
        keyboardVerticalOffset={0}
      >
        <BottomSheetScrollView
          style={{
            flex: 1,
            paddingHorizontal: screenPaddingHorizontal,
            flexDirection: "column",
          }}
        >
          <MapboxMap
            radius={radius}
            location={location}
            spaceName={name || undefined}
            color={colors[selectedColorIndex]}
            style={{
              marginBottom: template.margins.md,
            }}
          />
          <RadiusSection
            radius={radius}
            setRadius={handleRadiusChange}
            color={colors[selectedColorIndex]}
          />
          <NameSection
            spaceName={name}
            setSpaceName={handleNameChange}
            errors={spaceNameErrors}
          />
          <ColorSection
            selectedColorIndex={selectedColorIndex}
            setSelectedColorIndex={handleSelectedColorIndexChange}
          />
          <View
            style={{
              alignItems: "center",
              marginTop: template.margins.md + 10,
              marginBottom: insets.bottom + 20,
            }}
          >
            <PrimaryButton
              color={colors[selectedColorIndex]}
              onPress={handleSubmit}
              isDisabled={spaceNameErrors.length > 0}
            >
              Create Space
            </PrimaryButton>
          </View>
        </BottomSheetScrollView>
      </KeyboardAvoidingView>
    </View>
  );
};

const RadiusSection: FC<{
  radius: number;
  setRadius: (newRadius: number) => void;
  color: string;
}> = ({ setRadius, radius, color }) => {
  return (
    <View style={{ marginBottom: template.margins.md }}>
      <Label style={{ marginBottom: 10 }}>Radius</Label>
      <Slider
        initialValue={radius}
        maximumValue={100}
        onValueChange={setRadius}
        style={{ height: 35 }}
        thumbStyle={{ width: 35, backgroundColor: color }}
        trackStyle={{ height: 7, borderRadius: 4 }}
        minimumTrackTintColor={color}
        minimumValue={10}
        scaleFactor={1.7}
      />
    </View>
  );
};

const ColorSection: FC<{
  selectedColorIndex: number;
  setSelectedColorIndex: (newColorIndex: number) => void;
}> = ({ selectedColorIndex, setSelectedColorIndex }) => {
  return (
    <View>
      <Label style={{ marginBottom: 10 }}>Color</Label>
      <ColorPicker
        colors={colors}
        selectedIndex={selectedColorIndex}
        setSelectedColorIndex={setSelectedColorIndex}
        gapSize={gapSize}
        numberOfColumns={numberOfColumns}
        screenPaddingHorizontal={screenPaddingHorizontal}
      />
    </View>
  );
};

const NameSection: FC<{
  spaceName: string;
  setSpaceName: (newSpaceName: string) => void;
  errors: string[];
}> = ({ spaceName, setSpaceName, errors }) => {
  return (
    <View style={{ marginBottom: template.margins.md }}>
      <Label style={{ marginBottom: 10 }}>Name</Label>
      <TextInput
        errors={errors}
        text={spaceName}
        setText={setSpaceName}
        placeholder="Space Name"
      />
    </View>
  );
};
