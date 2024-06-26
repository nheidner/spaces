import { FC, useEffect, useState } from "react";
import { Button, Text, TextInput, View } from "react-native";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import { maybeCompleteAuthSession } from "expo-web-browser";
import { useAuthRequest } from "expo-auth-session/providers/google";
import {
  AppleAuthenticationButton,
  AppleAuthenticationButtonStyle,
  AppleAuthenticationButtonType,
  AppleAuthenticationScope,
  signInAsync,
} from "expo-apple-authentication";
import { auth } from "../../firebase";
import {
  GoogleAuthProvider,
  onAuthStateChanged,
  signInWithCredential,
  signOut,
  createUserWithEmailAndPassword,
  signInWithEmailAndPassword,
  sendEmailVerification,
  deleteUser,
  User,
  OAuthProvider,
} from "firebase/auth";
import { fetchApi } from "../utils/fetch_api";

maybeCompleteAuthSession();

export const Signin: FC<{}> = ({}) => {
  const [_, response, promptAsync] = useAuthRequest({
    iosClientId:
      "717063629448-ctoeh0a3vdaknng6cvmb2d23v1mjttk1.apps.googleusercontent.com",
    androidClientId:
      "717063629448-u808b96qbshbccogqoq7fnvf86fv33ne.apps.googleusercontent.com",
  });

  useEffect(() => {
    if (response?.type === "success") {
      const { id_token: idToken } = response.params;
      const credential = GoogleAuthProvider.credential(idToken);

      signInWithCredential(auth, credential)
        .then((userCredential) => {
          const firebaseIdToken = (userCredential as any)?._tokenResponse
            ?.idToken as string;

          if (!firebaseIdToken) {
            throw new Error("firebaseIdToken is undefined");
          }

          return fetchApi("/users", {
            method: "POST",
            body: JSON.stringify({
              idToken: firebaseIdToken,
            }),
          })
            .then((user) => {
              console.log("user :>> ", user);
            })
            .catch((error) => {
              console.error("error :>>", error);

              return signOut(auth);
            });
        })
        .catch((error) => console.error("error :>>", error));
    }

    if (response?.type === "error") {
      console.error("Authentication error:", response.error);
    }
  }, [response]);

  const insets = useSafeAreaInsets();

  return (
    <View style={{ paddingTop: insets.top, flex: 1, alignItems: "center" }}>
      <View style={{ marginTop: 25 }}>
        <Button
          title="Log in with Google"
          onPress={() => {
            promptAsync().catch((error) => console.error("error :>>", error));
          }}
        />
      </View>
      <AppleSignIn />
      <EmailSignIn />
      <EmailSignUp />
      <View style={{ marginTop: 25 }}>
        <Button
          title="sign out"
          onPress={() => {
            signOut(auth).catch((error) => console.error("error :>>", error));
          }}
        />
      </View>
    </View>
  );
};

const AppleSignIn: FC<{}> = () => {
  const handleSignIn = async () => {
    try {
      const credential = await signInAsync({
        requestedScopes: [
          AppleAuthenticationScope.FULL_NAME,
          AppleAuthenticationScope.EMAIL,
        ],
      });

      // if (!credential.identityToken) {
      //   throw new Error("Apple Sign-In failed - no identify token returned");
      // }

      // const provider = new OAuthProvider("apple.com");
      // const firebaseCredential = provider.credential({
      //   idToken: credential.identityToken,
      // });
      // const userCredential = await signInWithCredential(
      //   auth,
      //   firebaseCredential
      // );

      // const firebaseIdToken = (userCredential as any)?._tokenResponse
      //   ?.idToken as string;

      // if (!firebaseIdToken) {
      //   throw new Error("firebaseIdToken is undefined");
      // }

      // const user = await fetchApi<any>("/users", {
      //   method: "POST",
      //   body: JSON.stringify({
      //     idToken: firebaseIdToken,
      //   }),
      // });
    } catch (error: any) {
      console.log("error :>> ", error);
      if (error.code === "ERR_REQUEST_CANCELED") {
        // handle that the user canceled the sign-in flow
      } else {
        // handle other errors
      }
    }
  };

  return (
    <AppleAuthenticationButton
      buttonType={AppleAuthenticationButtonType.SIGN_IN}
      buttonStyle={AppleAuthenticationButtonStyle.BLACK}
      cornerRadius={10}
      style={{ height: 50, width: "60%", marginTop: 25 }}
      onPress={handleSignIn}
    />
  );
};
const EmailSignIn: FC<{}> = ({}) => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleSignIn = () => {
    signInWithEmailAndPassword(auth, email, password)
      .then((userCredential) => {
        console.log("userCredential.user :>> ", userCredential.user);
      })
      .catch((error) => {
        console.log("error :>> ", error);
      });
  };

  return (
    <View style={{ marginTop: 30 }}>
      <TextInput
        autoCapitalize="none"
        onChangeText={(newText) => setEmail(newText)}
        style={{
          borderWidth: 1,
          height: 40,
          width: 250,
          marginBottom: 10,
          padding: 5,
          fontSize: 20,
        }}
        placeholder="email"
      />
      <TextInput
        autoCapitalize="none"
        onChangeText={(newPassword) => setPassword(newPassword)}
        style={{
          borderWidth: 1,
          height: 40,
          width: 250,
          padding: 5,
          fontSize: 20,
        }}
        placeholder="password"
      />
      <Button title="sign in with email" onPress={handleSignIn} />
    </View>
  );
};

const EmailSignUp: FC<{}> = ({}) => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleSignUp = () => {
    createUserWithEmailAndPassword(auth, email, password)
      .then((userCredential) => {
        // Signed in
        const user = userCredential.user;
        const firebaseIdToken = (userCredential as any)?._tokenResponse
          ?.idToken as string;

        if (!firebaseIdToken) {
          throw new Error("firebaseIdToken is undefined");
        }

        // email verification
        sendEmailVerification(user).then(() => {
          console.log("sent email verification");
        });

        return fetchApi("/users", {
          method: "POST",
          body: JSON.stringify({
            idToken: firebaseIdToken,
          }),
        })
          .then((user) => {
            console.log("user :>> ", user);
          })
          .catch(() => deleteUser(auth.currentUser as User));
      })
      .catch((error) => {
        console.error("error :>> ", error);
      });
  };

  return (
    <View style={{ marginTop: 30 }}>
      <TextInput
        onChangeText={(newText) => setEmail(newText)}
        style={{
          borderWidth: 1,
          height: 40,
          width: 250,
          marginBottom: 10,
          padding: 5,
          fontSize: 20,
        }}
        placeholder="email"
      />
      <TextInput
        onChangeText={(newPassword) => setPassword(newPassword)}
        style={{
          borderWidth: 1,
          height: 40,
          width: 250,
          padding: 5,
          fontSize: 20,
        }}
        placeholder="password"
      />
      <Button title="sign up with email" onPress={handleSignUp} />
    </View>
  );
};
