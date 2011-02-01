package com.mfk.web.maker.client;


import com.google.gwt.event.dom.client.ClickEvent;
import com.google.gwt.event.dom.client.ClickHandler;
import com.google.gwt.user.client.ui.Button;
import com.google.gwt.user.client.ui.HTML;
import com.google.gwt.user.client.ui.VerticalPanel;


public class MfkPanel extends VerticalPanel {
  // The user-visible title for the entity
  public String title = "";
  // The user-visible image for the entity.
  public SearchImage image = new SearchImage();

  public MfkPanel(String title, SearchImage image) {
    this.setTitle(title);
    this.setImage(image);
    System.out.println("MfkPanel: title:" + title);
    this.addStyleName("mfkpanel");

  }

  public void setImage(SearchImage image) {
    this.image.setUrlAndQuery(image.getUrl(), image.getQuery());
    this.refresh();
  }

  public void setTitle(String title) {
    this.title = title;
    this.refresh();
  }

  /**
   * Refresh the UI elements of the page.
   */
  private void refresh() {
    this.clear();
    Button editButton = new Button("Edit");
    final MfkPanel outerThis = this;
    editButton.addClickHandler(new ClickHandler() {
      public void onClick(ClickEvent event) {
        MfkMaker.editDialog.editItem(outerThis);
      }
    });
    HTML title = new HTML(this.title);
    title.addStyleName("itemtitle");
    this.add(title);
    this.add(this.image);
    this.add(editButton);
  }

  public String toString() {
    return "<MfkPanel: " + this.title + ", url=" + this.image.getUrl()
        + ">";
  }

  /***
   * Tests if this MfkPanel is "too similar" to another.
   * Used for client-side validation.
   */
  public boolean similarTo(MfkPanel o) {
    if (o == null)
      return false;
    // Both title and URL must be distinct
    return this.title.equals(o.title) ||
        this.image.getUrl().equals(o.image.getUrl());
  }
}
