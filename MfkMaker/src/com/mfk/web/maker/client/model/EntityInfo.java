package com.mfk.web.maker.client.model;

public class EntityInfo {

  private final String name;
  private final String imageUrl;

  public EntityInfo(String name, String imageUrl) {
    this.name = name;
    this.imageUrl = imageUrl;
  }

  public String getName() {
    return name;
  }

  public String getImageUrl() {
    return imageUrl;
  }
}
