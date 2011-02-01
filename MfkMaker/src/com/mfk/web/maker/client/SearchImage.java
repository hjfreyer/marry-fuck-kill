package com.mfk.web.maker.client;


import com.google.gwt.user.client.ui.FlowPanel;
import com.google.gwt.user.client.ui.Image;


/**
 * A simple image-holder with constant width and height, designed specifically
 * to hold the images from an image search.
 */
public class SearchImage extends FlowPanel {
  private Image image;
  private String query;

  public SearchImage(String url, String query) {
    System.out.println("New SearchImage: " + url + ", " + query);
    this.image = new Image(url);
    this.query = new String(query);
    this.add(this.image);
    this.autoSize();
  }

  public SearchImage() {
    this.image = new Image();
    this.add(this.image);
    this.autoSize();
  }

  public void setUrlAndQuery(String url, String query) {
    System.out.println("SearchImage.setUrl: url=" + url + ", q=" + query);
    this.image.setUrl(url);
    this.query = new String(query);
  }

  public String getUrl() {
    return this.image.getUrl();
  }

  public String getQuery() {
    return this.query;
  }

  private void autoSize() {
    this.setWidth("145px");
    this.setHeight("145px");
    this.addStyleName("searchimage");
  }

  public String toString() {
    return "<SearchImage url=" + this.image.getUrl() + ", q=" + this.query
        + ">";
  }
}
