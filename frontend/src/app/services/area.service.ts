import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';

import { Area } from '../shared/models/area';

@Injectable()
export class AreaService {
  private _areas: Area[] = [];
  private _subject: BehaviorSubject<Area[]>;
  constructor() {
    if (window.localStorage['area']) {
      this._areas = <Area[]>JSON.parse(window.localStorage['area']);
    }
    this._subject = new BehaviorSubject(this._areas);
    this._subject.subscribe((areas: Area[]) => {
      window.localStorage['area'] = JSON.stringify(areas);
    });
  }
  getList(): Observable<Area[]> {
    return this._subject.asObservable();
  }
  save(areas: Area[]) {
    areas = areas.filter(area => {
      return (area.label && area.city);
    }).map((area: Area, index: number) => {
      if (!area.id) {
        area.id = `${Date.now()}-${index}`;
      }
      return area;
    });
    this._areas = [].concat(areas);
    this._subject.next(this._areas);
  }
  delete(id: string) {
    const result = this._areas.findIndex((area: Area) => {
      return (area.id === id);
    });
    if (result !== -1) {
      this._areas.splice(result, 1);
      this._subject.next(this._areas);
    }
  }
}