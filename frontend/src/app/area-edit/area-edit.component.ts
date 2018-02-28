import { Component, OnInit } from '@angular/core';
import { AreaService } from '../services/area.service';
import { Area } from '../shared/models/area';
import {
  transition,
  trigger,
  useAnimation
} from '@angular/animations';
import { slideFadeIn, slideFadeOut } from '../app.animations';

@Component({
  selector: 'app-area-edit',
  templateUrl: './area-edit.component.html',
  styleUrls: ['./area-edit.component.scss'],
  animations: [
    trigger('slideFade', [
      transition(':enter', [
        useAnimation(slideFadeIn)
      ]),
      transition(':leave', [
        useAnimation(slideFadeOut)
      ])
    ])
  ]
})
export class AreaEditComponent implements OnInit {
  public areas: Area[] = [];
  constructor(
    private areaService: AreaService
  ) { }

  ngOnInit() {
    this.areaService.getList().subscribe((areas: Area[]) => {
      this.areas = areas.slice();
    });
  }
  save(event) {
    event.preventDefault();
    this.areaService.save(this.areas);
  }
  delete(area: Area, index: number) {
    if (window.confirm(`${area.label} - 削除してもよろしいですか？`)) {
      if (!area.id) {
        this.areas.splice(index, 1);
      } else {
        this.areaService.delete(area.id);
      }
    }
  }
  addArea() {
    this.areas.push({
      id: null,
      label: '',
      city: ''
    });
  }
}