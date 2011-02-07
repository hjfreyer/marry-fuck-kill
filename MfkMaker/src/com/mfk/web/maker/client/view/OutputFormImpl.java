package com.mfk.web.maker.client.view;

import com.google.gwt.event.dom.client.HasClickHandlers;
import com.google.gwt.user.client.ui.Button;
import com.google.gwt.user.client.ui.FlowPanel;
import com.google.gwt.user.client.ui.FormPanel;
import com.google.gwt.user.client.ui.Hidden;
import com.google.gwt.user.client.ui.Panel;
import com.mfk.web.maker.client.model.EntityInfo;

public class OutputFormImpl implements OutputForm {

  // "null" means reload the page.
  private final FormPanel form = new FormPanel((String) null);
  private final Panel subPanel = new FlowPanel();

  private final Hidden n1 = new Hidden("n1");
  private final Hidden u1 = new Hidden("u1");
  private final Hidden q1 = new Hidden("q1");

  private final Hidden n2 = new Hidden("n2");
  private final Hidden u2 = new Hidden("u2");
  private final Hidden q2 = new Hidden("q2");

  private final Hidden n3 = new Hidden("n3");
  private final Hidden u3 = new Hidden("u3");
  private final Hidden q3 = new Hidden("q3");

  private final Button createButton = new Button("Create");

  public OutputFormImpl(Panel createButtonPanel, Panel formPanel) {
    formPanel.add(form);
    form.setWidget(subPanel);

    form.setAction("/make.do");
    form.setMethod(FormPanel.METHOD_POST);

    subPanel.add(n1);
    subPanel.add(u1);
    subPanel.add(q1);
    subPanel.add(n2);
    subPanel.add(u2);
    subPanel.add(q2);
    subPanel.add(n3);
    subPanel.add(u3);
    subPanel.add(q3);

    createButton.removeStyleName("gwt-Button");
    createButton.addStyleName("button");

    createButtonPanel.add(createButton);
  }

  @Override
  public void setEntity1(EntityInfo info) {
    n1.setValue(info.name);
    u1.setValue(info.imageUrl);
    q1.setValue(info.name);
  }

  @Override
  public void setEntity2(EntityInfo info) {
    n2.setValue(info.name);
    u2.setValue(info.imageUrl);
    q2.setValue(info.name);
  }

  @Override
  public void setEntity3(EntityInfo info) {
    n3.setValue(info.name);
    u3.setValue(info.imageUrl);
    q3.setValue(info.name);
  }

  @Override
  public void submit() {
    form.submit();
  }

  @Override
  public HasClickHandlers getSubmitButton() {
    return createButton;
  }

  @Override
  public void setClickable(boolean clickable) {
    if (clickable) {
      createButton.addStyleName("clickable");
    } else {
      createButton.removeStyleName("clickable");
    }
  }
}
