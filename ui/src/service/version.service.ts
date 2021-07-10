// rxjs
import { Observable } from "rxjs";

// memberdashboard
import { VersionResponse } from "../components/shared/types";
import { ENV } from "../env";
import { HTTPService } from "./http.service";

export class VersionService extends HTTPService {
  private readonly versionUrlSegment: string = ENV.api + "/version";

  getVersion(): Observable<VersionResponse> {
    return this.get<VersionResponse>(this.versionUrlSegment);
  }
}
