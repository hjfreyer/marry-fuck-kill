package com.mfk.web.maker.client;


import com.google.gwt.event.dom.client.ClickEvent;
import com.google.gwt.event.dom.client.ClickHandler;
import com.google.gwt.event.dom.client.KeyPressEvent;
import com.google.gwt.event.dom.client.KeyPressHandler;
import com.google.gwt.user.client.Timer;
import com.google.gwt.user.client.ui.Button;
import com.google.gwt.user.client.ui.DialogBox;
import com.google.gwt.user.client.ui.HTML;
import com.google.gwt.user.client.ui.HorizontalPanel;
import com.google.gwt.user.client.ui.Image;
import com.google.gwt.user.client.ui.TextBox;
import com.google.gwt.user.client.ui.VerticalPanel;


public class EditDialog extends DialogBox {
  private static String THROBBER_URL = "/s/loading.gif";
  private MfkPanel item = null;
  private SearchImage editImage = new SearchImage();
  private TextBox editTitle = new TextBox();
  private Image autoThrobber = new Image(EditDialog.THROBBER_URL);
  private HorizontalPanel searchThrobber = new HorizontalPanel();

  private Button searchButton = new Button("Search");
  private TextBox searchBox = new TextBox();

  // These are all bookkeeping for auto-search:

  // Last time we sent a search. Maintained by maybeSearch.
  private long lastSearchMillis = 0;
  // Last time the text box changed. Maintained by repeatingTimer.
  private long lastChangeMillis = 0;
  // Last search text. Maintained by maybeSearch.
  // NOTE: This is distinct from MfkMaker.resultSearchQuery -- this is used
  // only for maybeSearch's retry logic.
  private String lastSearch = "";

  // Timer that drives actual searching.
  private Timer repeatingTimer;

  // The expanding search panel.
  private VerticalPanel search = new VerticalPanel();

  public EditDialog() {
    this(false);
  }

  public EditDialog(boolean b) {
    super(b);

    HorizontalPanel searchControls = new HorizontalPanel();

    this.searchButton.addClickHandler(new ClickHandler() {
      public void onClick(ClickEvent event) {
        doSearch();
      }
    });
    this.searchBox.addKeyPressHandler(new KeyPressHandler() {
      public void onKeyPress(KeyPressEvent event) {
        if (event.getCharCode() == '\n' || event.getCharCode() == '\r') {
          doSearch();
        }
      }
    });

    searchControls.add(this.searchBox);
    searchControls.add(this.searchButton);

    HorizontalPanel moreImagesTitle = new HorizontalPanel();
    moreImagesTitle.add(new HTML("Search for more images:"));
    HTML hideLink = new HTML("[hide]");
    hideLink.addClickHandler(new ClickHandler() {
      public void onClick(ClickEvent event) {
        search.setVisible(false);
      }
    });
    hideLink.setStylePrimaryName("fakelink");
    moreImagesTitle.add(hideLink);

    this.search.add(new HTML("<hr>"));
    this.search.add(moreImagesTitle);
    this.search.add(searchControls);

    this.searchThrobber.add(new Image(EditDialog.THROBBER_URL));
    this.searchThrobber.add(new HTML("Loading..."));
    this.searchThrobber.setSpacing(10);
    this.searchThrobber.setVisible(false);
    this.search.add(searchThrobber);

    this.search.add(MfkMaker.resultPanel);
    this.search.setSpacing(5);
  }

  public void editItem(final MfkPanel item) {
    this.item = item;
    System.out.println("Showing dialog for :" + item);
    this.editImage.setUrlAndQuery(item.image.getUrl(), item.image
        .getQuery());
    this.editTitle.setText(item.title);
    this.autoThrobber.setVisible(false);
    this.search.setVisible(false);

    long now = System.currentTimeMillis();
    this.lastSearch = item.title;
    this.lastChangeMillis = this.lastSearchMillis = now;

    // This just keeps track of when the last change in the box was.
    // If it misses a keystroke, our time is a little old, but that's okay.
    // (We still throttle searches.)
    this.editTitle.addKeyPressHandler(new KeyPressHandler() {
      public void onKeyPress(KeyPressEvent event) {
        lastChangeMillis = System.currentTimeMillis();
        setAutoThrobber(true);
      }
    });

    // This checks repeatedly if we should search.
    this.repeatingTimer = new Timer() {
      public void run() {
        maybeSearch();
      }
    };
    this.repeatingTimer.scheduleRepeating(250);

    VerticalPanel p = new VerticalPanel();
    p.setSpacing(5);
    // TODO: put this in CSS, come up with a well-reasoned value
    p.setWidth("600px");

    Button create = new Button("<b>Save</b>");
    create.addClickHandler(new ClickHandler() {
      public void onClick(ClickEvent e) {
        System.out.println("Should create item here");
        hide();
        // update the existing item
        item.setTitle(editTitle.getText());
        item.setImage(editImage);
      }
    });

    Button cancel = new Button("Cancel");
    cancel.addClickHandler(new ClickHandler() {
      public void onClick(ClickEvent e) {
        hide();
      }
    });

    p.add(new HTML("<b>Name:</b>"));
    HorizontalPanel titlePanel = new HorizontalPanel();
    titlePanel.add(editTitle);
    titlePanel.add(autoThrobber);
    p.add(titlePanel);
    p.add(new HTML("<b>Image:</b>"));
    p.add(editImage);
    p.add(new HTML("Not the image you wanted?"));
    HTML link = new HTML("See more images.");
    link.addClickHandler(new ClickHandler() {
      public void onClick(ClickEvent event) {
        search.setVisible(!search.isVisible());
        if (search.isVisible()) {
          searchBox.setText(editTitle.getText());
          doSearch();
        }
      }
    });
    link.setStylePrimaryName("fakelink");
    p.add(link);

    HorizontalPanel buttonPanel = new HorizontalPanel();
    buttonPanel.setSpacing(5);
    buttonPanel.add(create);
    buttonPanel.add(cancel);
    p.add(buttonPanel);
    p.add(search);

    this.setWidget(p);
    this.show();
    this.center();
    this.editTitle.setFocus(true);
  }

  /**
   * This is the logic that determines if we should search.
   */
  public void maybeSearch() {
    long now = System.currentTimeMillis();
    String text = this.editTitle.getText();

    if (now - this.lastChangeMillis > 250) {
      if (now - this.lastSearchMillis > 1000) {
        this.lastSearchMillis = now;
        if (!text.equals(this.lastSearch) && !text.isEmpty()) {
          this.lastSearch = text;
          System.out.println("maybeSearch: <" + text + ">");
          this.doAutoSearch();
        }
      }
    }
  }

  private void doAutoSearch() {
    String text = this.editTitle.getText();
    this.searchBox.setText(text);
    this.doSearch();
  }

  private void doSearch() {
    MfkMaker.editDialog.setSearchThrobber(true);
    MfkMaker.resultPanel.clear();
    MfkMaker.results.clear();
    MfkMaker.resultsSearchQuery = this.searchBox.getText();
    MfkMaker.imageSearch.execute(this.searchBox.getText());
  }

  /**
   * Turn on or off the throbber.
   *
   * @param enabled
   */
  public void setAutoThrobber(boolean enabled) {
    this.autoThrobber.setVisible(enabled);
  }

  public void setSearchThrobber(boolean enabled) {
    this.searchThrobber.setVisible(enabled);
  }

  /**
   * Set the image for the item under edit.
   */
  public void setImage(Image image) {
    System.out.println("Set edit image url = " + image.getUrl()
        + " (from = " + MfkMaker.resultsSearchQuery + ")");
    this.editImage.setUrlAndQuery(image.getUrl(),
        MfkMaker.resultsSearchQuery);
  }

  public void autoSetImage(Image image) {
    if (this.shouldAutoSet()) {
      this.setImage(image);
    }
  }

  private boolean shouldAutoSet() {
    return !this.search.isVisible();
  }

  public void hide() {
    super.hide();
    this.item = null;
    if (this.repeatingTimer != null)
      this.repeatingTimer.cancel();
  }

  public MfkPanel getItem() {
    return this.item;
  }
}
