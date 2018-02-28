import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PugeventComponent } from './pugevent.component';

describe('PugeventComponent', () => {
  let component: PugeventComponent;
  let fixture: ComponentFixture<PugeventComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PugeventComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PugeventComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
