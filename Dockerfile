FROM scratch
ADD  eventscollector /
EXPOSE 13000
CMD ["/eventscollector"]