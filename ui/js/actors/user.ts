function URI(endpoint: string): string {
  return `${endpoint}`;
}

export class UserActor {
  userprofile: UserActor.UserProfile | null = null;
  async message(
    message: UserActor.MessageTypes,
    options: UserActor.RegisterRequest | UserActor.LoginRequest
  ): Promise<void | Boolean | UserActor.UserProfile | null> {
    switch (message) {
      case UserActor.MessageTypes.RegisterUser:
        const opts = options as UserActor.RegisterRequest;
        return await this.registerUser(
          opts.username,
          opts.password,
          opts.email
        );
      case UserActor.MessageTypes.Login:
        return await this.login(options.username, options.password);
      case UserActor.MessageTypes.GetUser:
        return await this.getUser();
      default:
        console.error("invalid message type sent to UserActor");
        return;
    }
  }

  resetUserProfile() {
    this.userprofile = null;
    return this.userprofile;
  }

  async getUser(): Promise<UserActor.UserProfile | null> {
    if (this.userprofile) return this.userprofile;

    const userResponse: Response | void = await fetch(URI("/api/user"));
    if (userResponse.status != 200) return this.resetUserProfile();
    if (userResponse.redirected) return this.resetUserProfile();

    const json: void | UserActor.UserProfile = await userResponse.json();

    if (!json) {
      this.userprofile = null;
      return this.userprofile;
    }
    this.userprofile = {
      username: json.username,
      email: json.email,
    };
    return await this.userprofile;
  }

  async login(username: string, password: string): Promise<Boolean> {
    const opts: UserActor.LoginRequest = {
      username: username,
      password: password,
    };

    const loginResponse: Response | void = await fetch(URI("/signin"), {
      method: "POST",
      body: JSON.stringify(opts),
    });
    if (loginResponse.status == 200) {
      return true;
    }

    return false;
  }

  async registerUser(
    username: string,
    password: string,
    email: string
  ): Promise<Boolean> {
    const opts: UserActor.RegisterRequest = {
      username: username,
      password: password,
      email: email,
    };
    const registerResponse: Response | void = await fetch(URI("/register"), {
      method: "POST",
      body: JSON.stringify(opts),
    });

    if (registerResponse.status == 200) {
      return true;
    }

    return false;
  }
}

export namespace UserActor {
  export enum MessageTypes {
    RegisterUser,
    Login,
    GetUser,
  }
  export interface RegisterRequest {
    username: string;
    password: string;
    email: string;
  }

  export interface LoginRequest {
    username: string;
    password: string;
  }

  export interface UserProfile {
    username: String;
    email: String;
  }
}
