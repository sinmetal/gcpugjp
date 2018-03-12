import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs/Observable";
import {OrganizationListAPIResponse} from "../shared/models/organization";

@Injectable()
export class OrganizationService {
  private API = '/api/1';

  constructor(public http: HttpClient) { }

  list(): Observable<OrganizationListAPIResponse> {
    return this.http.get<OrganizationListAPIResponse>(`${this.API}/organization`, {});
  }

}
