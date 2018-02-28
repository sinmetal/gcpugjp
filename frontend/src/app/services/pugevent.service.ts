import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/map';

import {PugeventListAPIResponse} from '../shared/models/Pugevent';

@Injectable()
export class PugeventService {
  private API = '/api/1';

  constructor(public http: HttpClient) {}

  list(): Observable<PugeventListAPIResponse> {
    return this.http.get<PugeventListAPIResponse>(`${this.API}/event`, {});
  }
}
